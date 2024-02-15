package protogen

import (
	"fmt"
	"strings"

	"github.com/gookit/color"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/util/stringx"
)

// BuildRPCs gen rpcs to string
func BuildRPCs(api *spec.ApiSpec, apiName string) (string, bool) {
	var builder strings.Builder
	var messageBuilder strings.Builder
	var hasEmpty bool
	methodMap := make(map[string]struct{})
	messageMap := make(map[string]struct{})

	serviceName := stringx.From(apiName).Title()
	builder.WriteString(fmt.Sprintf("// %s 服务\nservice %s {\n", serviceName, serviceName))
	for i, group := range api.Service.Groups {
		if i > 0 {
			builder.WriteByte('\n')
		}
		for _, route := range group.Routes {
			r, mf := parseRPC(route)
			if mf.MessageName != "" {
				if _, ok := messageMap[mf.MessageName]; !ok {
					messageBuilder.WriteString(fmt.Sprintf("%smessage %s {\n%srepeated %s %s = 1; // %s\n}\n\n",
						mf.MessageComment, mf.MessageName, indent, mf.FieldType, mf.FieldName, mf.Comment))
					messageMap[mf.MessageName] = struct{}{}
				}
			}
			if _, ok := methodMap[r.Method]; !ok {
				builder.WriteString(fmt.Sprintf("%s%s\n%srpc %s (%s) returns (%s);\n",
					indent, r.Doc, indent, r.Method, r.Request, r.Response))
				if r.HasEmpty {
					hasEmpty = true
				}
				methodMap[r.Method] = struct{}{}
			} else {
				color.Red.Printf("duplicate handler name, handler: %s, method: %s, path: %s, please rename it.\n",
					route.Handler, route.Method, route.Path)
			}
		}
	}
	builder.WriteByte('}')

	return messageBuilder.String() + builder.String(), hasEmpty
}

type rpc struct {
	Doc      string
	Method   string
	Request  string
	Response string
	HasEmpty bool
}

func parseRPC(route spec.Route) (rpc, messageField) {
	var mf messageField
	hasEmpty := false
	method := stringx.From(getHandlerBaseName(route)).Title()

	request := route.RequestTypeName()
	if request == "" {
		request = "Empty"
		hasEmpty = true
	}

	response := route.ResponseTypeName()
	if response == "" {
		response = "Empty"
		hasEmpty = true
	} else if strings.HasPrefix(response, "[]") {
		replacer := strings.NewReplacer("请求", "响应", "Req", "Resp")
		responseType := trimPrefix(response)
		mf = messageField{
			FieldName:      "results",
			FieldType:      responseType,
			Comment:        "结果",
			IsRepeated:     true,
			IsPointer:      false,
			MessageName:    responseType + "Resp",
			MessageComment: "// " + responseType + "Resp 详情信息\n",
		}
		if request != "Empty" {
			mf.MessageName = strings.TrimSuffix(request, "Req") + "Resp"
			if docs := route.RequestType.Documents(); len(docs) > 0 {
				mf.MessageComment = replacer.Replace(docs[len(docs)-1]) + "\n"
			}
		}
		response = mf.MessageName
	}

	docPrefix, docSuffix := "// "+method+" ", "方法"
	if route.AtDoc.Text != "" {
		docSuffix = strings.Trim(route.AtDoc.Text, `"`)
	} else if route.AtDoc.Properties != nil {
		docSuffix = strings.Trim(route.AtDoc.Properties["summary"], `"`)
	}
	if VarStringRemoveBeforeDelimiter != "" {
		if splits := strings.Split(docSuffix, VarStringRemoveBeforeDelimiter); len(splits) > 0 {
			docSuffix = splits[len(splits)-1]
		}
	}
	doc := docPrefix + docSuffix

	return rpc{
		Doc:      doc,
		Method:   method,
		Request:  request,
		Response: response,
		HasEmpty: hasEmpty,
	}, mf
}

func getHandlerBaseName(route spec.Route) string {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler
}
