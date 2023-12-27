{{.head}}

package server

import (
	{{if .notStream}}"context"{{end}}

	{{.imports}}
)

// {{.server}}Server {{.server}} 服务器结构
type {{.server}}Server struct {
	svcCtx *svc.ServiceContext
	{{.unimplementedServer}}
}

// New{{.server}}Server 新建 {{.server}} 服务器
func New{{.server}}Server(svcCtx *svc.ServiceContext) *{{.server}}Server {
	return &{{.server}}Server{
		svcCtx: svcCtx,
	}
}

{{.funcs}}
