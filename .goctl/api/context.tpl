package svc

import (
	{{.configImport}}
)

// ServiceContext 服务上下文
type ServiceContext struct {
	m      sync.RWMutex
	c      {{.config}}
	Apollo *apollo.Apollo
	{{.middleware}}
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c {{.config}}) *ServiceContext {
	a := apollo.MustNewApollo(&c.Apollo)
	InitConfig(&c, a)

	sc := &ServiceContext{
		c:      c,
		Apollo: a,
		{{.middlewareAssignment}}
	}

	return sc
}

// InitConfig 初始化网关全局配置
func InitConfig(c *config.Config, a *apollo.Apollo) {
	apollo.MustUnmarshalYaml(a.GetNamespaceContent("gateway.yaml"), &c)
}

// Config 获取全局配置
func (sc *ServiceContext) Config() config.Config {
	sc.m.RLock()
	defer sc.m.RUnlock()

	return sc.c
}

// OnUpdate 配置更新事件处理
func (sc *ServiceContext) OnUpdate(event *agollo.ChangeEvent) {
	sc.m.Lock()
	defer sc.m.Unlock()

	if event.Namespace != "gateway.yaml" {
		return
	}

	for _, change := range event.Changes {
		if change.ChangeType == agollo.MODIFY && change.NewValue != "" && change.Key == "content" {
			err := apollo.UnmarshalYaml(change.NewValue, &sc.c)
			if err != nil {
				logx.Errorf("OnUpdate handle change event err, err: %v, change: %s", err, change)
			} else {
				logx.Infof("OnUpdate receive change event, change: %s", change)
			}
		}
	}
}
