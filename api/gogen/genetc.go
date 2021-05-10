package gogen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sliveryou/goctl/api/spec"
	"github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/util/format"
)

const (
	defaultPort = 8888
	etcDir      = "etc"
	etcTemplate = `Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
`
)

func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)
	serviceName := service.Name
	if i := strings.Index(serviceName, "service"); i > 0 {
		serviceName = strings.TrimSuffix(serviceName[:i], "-")
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          etcDir,
		filename:        fmt.Sprintf("%s.yaml", filename),
		templateName:    "etcTemplate",
		category:        category,
		templateFile:    etcTemplateFile,
		builtinTemplate: etcTemplate,
		data: map[string]string{
			"serviceName": serviceName,
			"host":        host,
			"port":        port,
		},
	})
}
