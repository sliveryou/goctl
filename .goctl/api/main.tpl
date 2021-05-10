package main

import (
	"flag"
	"fmt"
	"os"

	{{.importPackages}}
)

var (
	// projectName 项目名称
	projectName = "Project Name"
	// projectVersion 项目版本
	projectVersion = "1.0.0"
	// goVersion go版本
	goVersion = ""
	// gitCommit git提交commit id
	gitCommit = ""
	// buildTime 编译时间
	buildTime = ""

	// configFile 配置文件路径
	configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")
	// isShowVersion 是否显示项目版本信息
	isShowVersion = flag.Bool("v", false, "show project version")
)

// @Title 网关接口
// @Version 1.0.0
// @Description 网关平台
// @SecurityDefinitions.ApiKey ApiKeyAuth
// @In header
// @Name Authorization
// @BasePath /
func main() {
	flag.Parse()
	showVersion(*isShowVersion)

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// showVersion 显示项目版本信息
func showVersion(isShow bool) {
	if isShow {
		fmt.Printf("Project: %s\n", projectName)
		fmt.Printf(" Version: %s\n", projectVersion)
		fmt.Printf(" Go Version: %s\n", goVersion)
		fmt.Printf(" Git Commit: %s\n", gitCommit)
		fmt.Printf(" Build Time: %s\n", buildTime)
		os.Exit(0)
	}
}
