
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}) {{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}}, error) {
	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx, in)
}
