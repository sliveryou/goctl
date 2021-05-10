{{.head}}

//go:generate mockgen -destination ./{{.name}}_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/zrpc"
)

type (
	{{.alias}}

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}
)

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

{{.functions}}
