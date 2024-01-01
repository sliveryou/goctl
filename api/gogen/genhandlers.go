package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/format"
	"github.com/sliveryou/goctl/util/pathx"
)

const defaultLogicPackage = "logic"

//go:embed handler.tpl
var handlerTemplate string

type handlerInfo struct {
	PkgName            string
	ImportPackages     string
	ImportHttpxPackage string
	HandlerName        string
	PathName           string
	MethodName         string
	Tag                string
	Summary            string
	ResponseType       string
	RequestType        string
	LogicName          string
	LogicType          string
	Call               string
	HasResp            bool
	HasRequest         bool
	HasSecurity        bool
	HasRequestBody     bool
	SwagParams         []swagParam
	ResponseParamType  string
	ResponseDataType   string
	Accept             string
}

func genHandler(dir, rootPkg, srvName string, cfg *config.Config, group spec.Group, route spec.Route) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = strings.Title(handler)
		logicName = pkgName
	}

	tag := "Tag"
	if a := group.GetAnnotation("tag"); a != "" {
		tag = strings.Trim(a, `"`)
	} else {
		tag = strings.ToLower(srvName) + "/" + group.GetAnnotation("group")
	}
	summary := "Summary"
	if route.AtDoc.Properties != nil {
		summary = strings.Trim(route.AtDoc.Properties["summary"], `"`)
	}
	hasSecurity := false
	if ja := group.GetAnnotation("jwt"); ja != "" {
		hasSecurity = true
	} else if ma := strings.ToLower(group.GetAnnotation("middleware")); ma != "" {
		if strings.Contains(ma, "jwt") || strings.Contains(ma, "auth") {
			hasSecurity = true
		}
	}

	reg := regexp.MustCompile(`/:([^/]+)`)
	pathName := reg.ReplaceAllString(strings.TrimSpace(route.Path), "/{${1}}")
	methodName := strings.ToLower(strings.TrimSpace(route.Method))

	sps, existJsonTag, existFormTag := getSwagParams(route)
	respType := util.Title(route.ResponseTypeName())
	respDataType := strings.TrimPrefix(strings.TrimPrefix(respType, "[]"), "*")
	respParamType := "{object}"
	if len(respType) != len(respDataType) {
		respParamType = "{array}"
	}
	accept := "application/json"
	if methodName != "get" && !existJsonTag && existFormTag {
		accept = "multipart/form-data"
	}

	return doGenToFile(dir, handler, cfg, group, route, handlerInfo{
		PkgName:           pkgName,
		ImportPackages:    genHandlerImports(group, route, rootPkg),
		HandlerName:       handler,
		PathName:          pathName,
		MethodName:        methodName,
		Tag:               tag,
		Summary:           summary,
		ResponseType:      util.Title(route.ResponseTypeName()),
		RequestType:       util.Title(route.RequestTypeName()),
		LogicName:         logicName,
		LogicType:         strings.Title(getLogicName(route)),
		Call:              strings.Title(strings.TrimSuffix(handler, "Handler")),
		HasResp:           len(route.ResponseTypeName()) > 0,
		HasRequest:        len(route.RequestTypeName()) > 0,
		HasSecurity:       hasSecurity,
		HasRequestBody:    len(route.RequestTypeName()) > 0 && existJsonTag,
		SwagParams:        sps,
		ResponseParamType: respParamType,
		ResponseDataType:  respDataType,
		Accept:            accept,
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
			if err := genHandler(dir, rootPkg, api.Service.Name, cfg, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}
	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
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

func getSwagParams(route spec.Route) ([]swagParam, bool, bool) {
	rt := route.RequestType
	method := strings.ToLower(strings.TrimSpace(route.Method))

	if rt != nil {
		if ds, ok := rt.(spec.DefineStruct); ok {
			return parseSwagParams(ds, method)
		}
	}

	return nil, false, false
}

func parseSwagParams(ds spec.DefineStruct, method string) ([]swagParam, bool, bool) {
	var sps []swagParam
	existJsonTag := false
	existFormTag := false

	for _, m := range ds.Members {
		switch mt := convertSpec(m.Type).(type) {
		case spec.PrimitiveType:
			tags, err := spec.Parse(m.Tag)
			if err != nil {
				fmt.Printf("request type: %s parse tags err, err: %v\n", ds.RawName, err)
				break
			}
			for _, tag := range tags.Tags() {
				paramType, findJsonTag, findFormTag := getParamType(tag.Key, method)
				if findJsonTag {
					existJsonTag = true
				}
				if findFormTag {
					existFormTag = true
				}
				if paramType != "" {
					sps = append(sps, swagParam{
						ParamName:   tag.Name,
						ParamType:   paramType,
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
					paramType, findJsonTag, findFormTag := getParamType(tag.Key, method)
					if findJsonTag {
						existJsonTag = true
					}
					if findFormTag {
						existFormTag = true
					}
					if paramType == "query" || paramType == "formData" {
						sps = append(sps, swagParam{
							ParamName:   tag.Name,
							ParamType:   paramType,
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
				s, jt, ft := parseSwagParams(mt, method)
				if jt {
					existJsonTag = true
				}
				if ft {
					existFormTag = true
				}
				sps = append(sps, s...)
			}
		}
	}

	return sps, existJsonTag, existFormTag
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

func convertSpec(t spec.Type) spec.Type {
	var tt spec.PointerType
	var ok = true

	for ok {
		tt, ok = t.(spec.PointerType)
		if ok {
			t = tt.Type
		}
	}

	return t
}

func getComment(comment string) string {
	return strings.TrimSpace(strings.TrimPrefix(comment, "//"))
}

func getParamType(key, method string) (paramType string, findJsonTag, findFormTag bool) {
	switch key {
	case "path":
		paramType = "path"
	case "form":
		if method == "get" {
			paramType = "query"
		} else {
			paramType = "formData"
		}
		findFormTag = true
	case "header":
		paramType = "header"
	case "json":
		findJsonTag = true
	}

	return
}
