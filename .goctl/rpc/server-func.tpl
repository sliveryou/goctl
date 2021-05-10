
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in {{.request}}) ({{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
