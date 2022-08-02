package protogen

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/sliveryou/goctl/api/spec"
)

// BuildRPCs gen rpcs to string
func BuildRPCs(api *spec.ApiSpec) (string, bool) {
	var builder strings.Builder
	var hasEmpty bool
	methodMap := make(map[string]struct{})

	builder.WriteString("// RPC 相关服务\nservice RPC {\n")
	for i, group := range api.Service.Groups {
		if i > 0 {
			builder.WriteByte('\n')
		}
		for _, route := range group.Routes {
			r := parseRPC(route)
			if _, ok := methodMap[r.Method]; !ok {
				builder.WriteString(fmt.Sprintf("%s%s\n%srpc %s (%s) returns (%s);\n",
					indent, r.Doc, indent, r.Method, r.Request, r.Response))
				if r.HasEmpty {
					hasEmpty = true
				}
				methodMap[r.Method] = struct{}{}
			} else {
				fmt.Println(aurora.Red(fmt.Sprintf("duplicate handler name, handler: %s, method: %s, path: %s, please rename it.",
					route.Handler, route.Method, route.Path)))
			}
		}
	}
	builder.WriteByte('}')

	return builder.String(), hasEmpty
}

type rpc struct {
	Doc      string
	Method   string
	Request  string
	Response string
	HasEmpty bool
}

func parseRPC(route spec.Route) rpc {
	hasEmpty := false
	method := strings.Title(getHandlerBaseName(route))
	request := route.RequestTypeName()
	if request == "" {
		request = "Empty"
		hasEmpty = true
	}
	response := route.ResponseTypeName()
	if response == "" {
		response = "Empty"
		hasEmpty = true
	}
	doc := "// " + method + " 方法"
	if route.AtDoc.Properties != nil {
		doc = "// " + method + " " + strings.Trim(route.AtDoc.Properties["summary"], `"`)
	}

	return rpc{
		Doc:      doc,
		Method:   method,
		Request:  request,
		Response: response,
		HasEmpty: hasEmpty,
	}
}

func getHandlerBaseName(route spec.Route) string {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler
}
