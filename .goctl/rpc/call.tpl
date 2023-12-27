{{.head}}

package {{.filePackage}}

import (
	"context"

	{{.pbPackage}}
	{{if ne .pbPackage .protoGoPackage}}{{.protoGoPackage}}{{end}}

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// 类型定义
type (
	{{.alias}}

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}
)

// New{{.serviceName}} 新建 {{.serviceName}} 客户端
func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

{{.functions}}
