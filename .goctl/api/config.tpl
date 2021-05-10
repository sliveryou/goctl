package config

import {{.authImport}}

// Config 全局相关配置
type Config struct {
	rest.RestConf
	Apollo apollo.Config
	{{.auth}}
}
