package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/rpc/parser"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/format"
	"github.com/sliveryou/goctl/util/stringx"
)

const svcTemplate = `package svc

import {{.imports}}

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:c,
	}
}
`

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *DefaultGenerator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetSvc()
	svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, svcFilename+".go")
	text, err := util.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	serviceName := strings.ToLower(stringx.From(ctx.GetServiceName().Source()).ToCamel())
	if i := strings.Index(serviceName, "service"); i > 0 {
		serviceName = strings.TrimSuffix(serviceName[:i], "-")
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"imports":     fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
		"serviceName": serviceName,
	}, fileName, false)
}
