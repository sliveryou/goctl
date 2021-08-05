{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context{{if .hasReq}},in *{{.pbRequest}}{{end}}) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error)