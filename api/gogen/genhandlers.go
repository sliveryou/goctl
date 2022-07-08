package gogen

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/internal/version"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/format"
	"github.com/sliveryou/goctl/vars"
)

const (
	defaultLogicPackage = "logic"
	handlerTemplate     = `package {{.PkgName}}

import (
	"net/http"

	{{if .After1_1_10}}"github.com/tal-tech/go-zero/rest/httpx"{{end}}
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
			httpx.Error(w, err)
		} else {
			{{if .HasResp}}httpx.OkJson(w, resp){{else}}httpx.Ok(w){{end}}
		}
	}
}
`
)

type handlerInfo struct {
	PkgName        string
	ImportPackages string
	HandlerName    string
	PathName       string
	MethodName     string
	Tag            string
	Summary        string
	ResponseType   string
	RequestType    string
	LogicName      string
	LogicType      string
	Call           string
	HasResp        bool
	HasRequest     bool
	HasSecurity    bool
	After1_1_10    bool
	HasRequestBody bool
	SwagParams     []swagParam
}

func genHandler(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = strings.Title(handler)
		logicName = pkgName
	}
	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}
	tag := "Tag"
	if a := group.GetAnnotation("tag"); a != "" {
		tag = strings.TrimSuffix(strings.TrimPrefix(a, "\""), "\"")
	}
	summary := "Summary"
	if route.AtDoc.Properties != nil {
		summary = strings.TrimSuffix(strings.TrimPrefix(route.AtDoc.Properties["summary"], "\""), "\"")
	}
	hasSecurity := false
	if ja := group.GetAnnotation("jwt"); ja != "" {
		hasSecurity = true
	} else if ma := strings.ToLower(group.GetAnnotation("middleware")); ma != "" {
		if strings.Contains(ma, "jwt") || strings.Contains(ma, "auth") {
			hasSecurity = true
		}
	}

	goctlVersion := version.GetGoctlVersion()
	// todo(anqiansong): This will be removed after a certain number of production versions of goctl (probably 5)
	after1_1_10 := version.IsVersionGreaterThan(goctlVersion, "1.1.10")

	reg := regexp.MustCompile(`/:([^/]+)`)
	pathName := reg.ReplaceAllString(strings.TrimSpace(route.Path), "/{${1}}")
	methodName := strings.ToLower(strings.TrimSpace(route.Method))
	sps := getSwagParams(route)

	return doGenToFile(dir, handler, cfg, group, route, handlerInfo{
		PkgName:        pkgName,
		ImportPackages: genHandlerImports(group, route, parentPkg),
		HandlerName:    handler,
		PathName:       pathName,
		MethodName:     methodName,
		Tag:            tag,
		Summary:        summary,
		ResponseType:   util.Title(route.ResponseTypeName()),
		RequestType:    util.Title(route.RequestTypeName()),
		LogicName:      logicName,
		LogicType:      strings.Title(getLogicName(route)),
		Call:           strings.Title(strings.TrimSuffix(handler, "Handler")),
		HasResp:        len(route.ResponseTypeName()) > 0,
		HasRequest:     len(route.RequestTypeName()) > 0,
		HasSecurity:    hasSecurity,
		After1_1_10:    after1_1_10,
		HasRequestBody: len(route.RequestTypeName()) > 0 && len(sps) == 0,
		SwagParams:     sps,
	})
}

func doGenToFile(dir, handler string, cfg *config.Config, group spec.Group,
	route spec.Route, handleObj handlerInfo) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data:            handleObj,
	})
}

func genHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, rootPkg, cfg, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"",
		util.JoinPackages(parentPkg, getLogicFolderPath(group, route))))
	imports = append(imports, fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", util.JoinPackages(parentPkg, typesDir)))
	}

	currentVersion := version.GetGoctlVersion()
	// todo(anqiansong): This will be removed after a certain number of production versions of goctl (probably 5)
	if !version.IsVersionGreaterThan(currentVersion, "1.1.10") {
		imports = append(imports, fmt.Sprintf("\"%s/rest/httpx\"", vars.ProjectOpenSourceURL))
	}

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}

	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getLogicName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Logic"
}

type swagParam struct {
	ParamName   string
	ParamType   string
	DataType    string
	IsMandatory string
	Comment     string
	Attribute   string
}

func getSwagParams(route spec.Route) []swagParam {
	rt := route.RequestType
	methodName := strings.ToLower(strings.TrimSpace(route.Method))

	if rt != nil && methodName == "get" {
		if ds, ok := rt.(spec.DefineStruct); ok {
			return parseSwagParams(ds)
		}
	}

	return nil
}

func parseSwagParams(ds spec.DefineStruct) []swagParam {
	var sps []swagParam

	for _, m := range ds.Members {
		switch mt := m.Type.(type) {
		case spec.PrimitiveType:
			tags, err := spec.Parse(m.Tag)
			if err != nil {
				fmt.Printf("request type: %s parse tags err, err: %v\n", ds.RawName, err)
				break
			}
			for _, tag := range tags.Tags() {
				if tag.Key == "form" {
					sps = append(sps, swagParam{
						ParamName:   tag.Name,
						ParamType:   "query",
						DataType:    getDataType(mt.RawName),
						IsMandatory: strconv.FormatBool(!stringx.Contains(tag.Options, "optional")),
						Comment:     getComment(m.Comment),
					})
					break
				}
			}
		case spec.ArrayType:
			if vt, ok := mt.Value.(spec.PrimitiveType); ok {
				tags, err := spec.Parse(m.Tag)
				if err != nil {
					fmt.Printf("request type: %s parse tags err, err: %v\n", ds.RawName, err)
					break
				}
				for _, tag := range tags.Tags() {
					if tag.Key == "form" {
						sps = append(sps, swagParam{
							ParamName:   tag.Name,
							ParamType:   "query",
							DataType:    "[]" + getDataType(vt.RawName),
							IsMandatory: strconv.FormatBool(!stringx.Contains(tag.Options, "optional")),
							Comment:     getComment(m.Comment),
							Attribute:   "collectionFormat(multi)",
						})
						break
					}
				}
			}
		case spec.DefineStruct:
			if m.IsInline {
				sps = append(sps, parseSwagParams(mt)...)
			}
		}
	}

	return sps
}

func getDataType(dataType string) string {
	switch dataType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64", "complex64", "complex128":
		return "number"
	case "bool":
		return "boolean"
	}

	return ""
}

func getComment(comment string) string {
	return strings.TrimSpace(strings.TrimPrefix(comment, "//"))
}
