package main

import (
	"github.com/sliveryou/goctl/cmd"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
