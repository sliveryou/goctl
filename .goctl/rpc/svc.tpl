package svc

import {{.imports}}

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:c,
	}
}
