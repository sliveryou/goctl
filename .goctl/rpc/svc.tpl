package svc

import {{.imports}}

// ServiceContext 服务上下文
type ServiceContext struct {
	m      sync.RWMutex
	c      config.Config
	Apollo *apollo.Apollo
	DB     *gorm.DB
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	a := apollo.MustNewApollo(&c.Apollo)
	InitConfig(&c, a)

	db := gdb.MustNewDB(&c.DB)
	err := InitModel(context.Background(), db)
	if err != nil {
		panic(err)
	}

	sc := &ServiceContext{
		c:      c,
		Apollo: a,
		DB:     db,
	}

	return sc
}

// InitConfig 初始化服务全局配置
func InitConfig(c *config.Config, a *apollo.Apollo) {
	apollo.MustUnmarshalYaml(a.GetNamespaceContent("service.yaml"), &c)
	apollo.MustUnmarshalYaml(a.GetNamespaceValue("application", "DB"), &c)
	apollo.MustUnmarshalYaml(a.GetNamespaceValue("application", "Etcd"), &c.RpcServerConf)
}

// InitModel 初始化服务模型信息
func InitModel(ctx context.Context, db *gorm.DB) error {
	err := db.Set(
		"gorm:table_options",
		"ENGINE=InnoDB AUTO_INCREMENT=1 CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci",
	).AutoMigrate()
	if err != nil {
		return err
	}

	return nil
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

	if event.Namespace != "service.yaml" {
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
