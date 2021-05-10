package handler

import (
	"net/http"

	{{.ImportPackages}}
	"gitlab.33.cn/proof/backend-micro/pkg/xhttp"
)

// @Tags {{.Tag}}
// @Summary {{.Summary}}
{{- if .HasSecurity}}
// @Security ApiKeyAuth{{end}}
// @Accept application/json
// @Produce application/json
{{- if .HasRequest}}
// @Param data body types.{{.RequestType}} true "{{.Summary}}"{{end}}
// @Success 200 {object} types.Response{data=types.{{.ResponseType}}} "{"code":200,"msg":"OK","data":{}}"
// @Router {{.PathName}} [{{.MethodName}}]
func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := xhttp.Parse(r, &req); err != nil {
			xhttp.Error(w, r, err)
			return
		}
		
		{{end -}}
		l := logic.New{{.LogicType}}(r.Context(), ctx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
			xhttp.Error(w, r, err)
		} else {
			{{if .HasResp}}xhttp.OkJson(w, r, resp){{else}}httpx.Ok(w){{end}}
		}
	}
}
