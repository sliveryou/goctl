package config

import (
 	"github.com/tal-tech/go-zero/zrpc"

    "gitlab.33.cn/proof/backend-micro/pkg/apollo"
 	"gitlab.33.cn/proof/backend-micro/pkg/gdb"
 )

// Config 全局相关配置
type Config struct {
	zrpc.RpcServerConf
	Apollo apollo.Config
	DB     gdb.Config
}
