package gogen

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/format"
	"github.com/sliveryou/goctl/vars"
)

const (
	routesFilename = "routes"
	routesTemplate = `// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	{{.importPackages}}
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	{{.routesAdditions}}
}
`
	routesAdditionTemplate = `
	engine.AddRoutes(
		{{.routes}} {{.jwt}}{{.signature}}
	)
`
)

var mapping = map[string]string{
	"delete": "http.MethodDelete",
	"get":    "http.MethodGet",
	"head":   "http.MethodHead",
	"post":   "http.MethodPost",
	"put":    "http.MethodPut",
	"patch":  "http.MethodPatch",
}

type (
	group struct {
		routes           []route
		jwtEnabled       bool
		signatureEnabled bool
		authName         string
		middlewares      []string
	}
	route struct {
		method  string
		path    string
		handler string
	}
)

func genRoutes(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	var builder strings.Builder
	groups, err := getRoutes(api)
	if err != nil {
		return err
	}

	gt := template.Must(template.New("groupTemplate").Parse(routesAdditionTemplate))
	for _, g := range groups {
		var gbuilder strings.Builder
		gbuilder.WriteString("[]rest.Route{")
		for _, r := range g.routes {
			fmt.Fprintf(&gbuilder, `
		{
			Method:  %s,
			Path:    "%s",
			Handler: %s,
		},`,
				r.method, r.path, r.handler)
		}

		var jwt string
		if g.jwtEnabled {
			jwt = fmt.Sprintf("\n rest.WithJwt(serverCtx.Config.%s.AccessSecret),", g.authName)
		}
		var signature string
		if g.signatureEnabled {
			signature = "\n rest.WithSignature(serverCtx.Config.Signature),"
		}

		var routes string
		if len(g.middlewares) > 0 {
			gbuilder.WriteString("\n}...,")
			params := g.middlewares
			for i := range params {
				params[i] = "serverCtx." + params[i]
			}
			middlewareStr := strings.Join(params, ", ")
			routes = fmt.Sprintf("rest.WithMiddlewares(\n[]rest.Middleware{ %s }, \n %s \n),",
				middlewareStr, strings.TrimSpace(gbuilder.String()))
		} else {
			gbuilder.WriteString("\n},")
			routes = strings.TrimSpace(gbuilder.String())
		}

		if err := gt.Execute(&builder, map[string]string{
			"routes":    routes,
			"jwt":       jwt,
			"signature": signature,
		}); err != nil {
			return err
		}
	}

	routeFilename, err := format.FileNamingFormat(cfg.NamingFormat, routesFilename)
	if err != nil {
		return err
	}
	routeFilename = routeFilename + ".go"

	filename := path.Join(dir, handlerDir, routeFilename)
	os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          handlerDir,
		filename:        routeFilename,
		templateName:    "routesTemplate",
		category:        "",
		templateFile:    "",
		builtinTemplate: routesTemplate,
		data: map[string]string{
			"importPackages":  genRouteImports(rootPkg, api),
			"routesAdditions": strings.TrimSpace(builder.String()),
		},
	})
}

func genRouteImports(parentPkg string, api *spec.ApiSpec) string {
	importSet := collection.NewSet()
	importSet.AddStr(fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			folder := route.GetAnnotation(groupProperty)
			if len(folder) == 0 {
				folder = group.GetAnnotation(groupProperty)
				if len(folder) == 0 {
					continue
				}
			}
			importSet.AddStr(fmt.Sprintf("%s \"%s\"", toPrefix(folder), util.JoinPackages(parentPkg, handlerDir, folder)))
		}
	}
	imports := importSet.KeysStr()
	sort.Strings(imports)
	projectSection := strings.Join(imports, "\n\t")
	depSection := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL)
	return fmt.Sprintf("%s\n\n\t%s", projectSection, depSection)
}

func getRoutes(api *spec.ApiSpec) ([]group, error) {
	var routes []group

	for _, g := range api.Service.Groups {
		var groupedRoutes group
		for _, r := range g.Routes {
			handler := getHandlerName(r)
			handler = handler + "(serverCtx)"
			folder := r.GetAnnotation(groupProperty)
			if len(folder) > 0 {
				handler = toPrefix(folder) + "." + strings.ToUpper(handler[:1]) + handler[1:]
			} else {
				folder = g.GetAnnotation(groupProperty)
				if len(folder) > 0 {
					handler = toPrefix(folder) + "." + strings.ToUpper(handler[:1]) + handler[1:]
				}
			}
			groupedRoutes.routes = append(groupedRoutes.routes, route{
				method:  mapping[r.Method],
				path:    r.Path,
				handler: handler,
			})
		}

		jwt := g.GetAnnotation("jwt")
		if len(jwt) > 0 {
			groupedRoutes.authName = jwt
			groupedRoutes.jwtEnabled = true
		}
		signature := g.GetAnnotation("signature")
		if signature == "true" {
			groupedRoutes.signatureEnabled = true
		}
		middleware := g.GetAnnotation("middleware")
		if len(middleware) > 0 {
			for _, item := range strings.Split(middleware, ",") {
				groupedRoutes.middlewares = append(groupedRoutes.middlewares, item)
			}
		}
		routes = append(routes, groupedRoutes)
	}

	return routes, nil
}

func toPrefix(folder string) string {
	return strings.ReplaceAll(folder, "/", "")
}
