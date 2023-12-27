{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}},error)