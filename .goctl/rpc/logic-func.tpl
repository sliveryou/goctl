{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} (in {{.request}}) ({{.response}}, error) {
	// todo: add your logic here and delete this line
	
	return &{{.responseType}}{}, nil
}
