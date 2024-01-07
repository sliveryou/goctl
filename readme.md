# goctl

[![Go Report Card](https://goreportcard.com/badge/github.com/sliveryou/goctl)](https://goreportcard.com/report/github.com/sliveryou/goctl)
[![goproxy](https://goproxy.cn/stats/github.com/sliveryou/goctl/badges/download-count.svg)](https://goproxy.cn/stats/github.com/sliveryou/goctl/badges/download-count.svg)
[![Release](https://img.shields.io/github/v/release/sliveryou/goctl.svg?style=flat-square)](https://github.com/sliveryou/goctl)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 重要！

原仓库地址：https://github.com/zeromicro/go-zero/tree/master/tools/goctl  
主要根据原仓库的 goctl 对 api 和 rpc 代码的生成等做了一些定制化的修改  
基于 goctl 1.6.1 版本：https://github.com/zeromicro/go-zero/tree/v1.6.1/tools/goctl  

## 改动

1. 优化：当有一个 api 文件，如 `base.api` 文件 import 了多个 api 文件时，这些 api 文件可以跨文件共享定义好的结构体，方便一些公有结构体放在 `common.api` 中，其他文件共用
2. 优化：增加 `goctl api proto` 命令，可以基于 api 文件生成 proto 文件，使用例子：`goctl api proto --api base.api --dir .`
3. 优化：支持生成 handler 相关文件时增加 swag 注解，需使用项目内模板 `.goctl`，使用例子：`goctl api go --api base.api --dir . --home .goctl`
   1. 配合 [swag](https://github.com/swaggo/swag) 工具可以生成 swagger 文档和 swagger 服务，使用例子：`swag init -d . -g main.go -p snakecase --ot go,json,yaml -o docs`
   2. 可以将项目内模板作为 goctl 执行默认模板，命令：`mkdir -p ~/.goctl/1.6.1 && cp -r .goctl/* ~/.goctl/1.6.1`
