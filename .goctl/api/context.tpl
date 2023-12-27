package svc

import (
	{{.configImport}}
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config {{.config}}
	{{.middleware}}
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
		{{.middlewareAssignment}}
	}
}
