package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"

	{{.imports}}
)

// {{.logicName}} {{.comment}}上下文
type {{.logicName}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// New{{.logicName}} 新建{{.comment}}上下文
func New{{.logicName}}(ctx context.Context,svcCtx *svc.ServiceContext) *{{.logicName}} {
	return &{{.logicName}}{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
{{.functions}}
