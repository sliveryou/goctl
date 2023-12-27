package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/sliveryou/goctl/config"
	"github.com/sliveryou/goctl/rpc/parser"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/format"
	"github.com/sliveryou/goctl/util/pathx"
	"github.com/sliveryou/goctl/util/stringx"
)

//go:embed svc.tpl
var svcTemplate string

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *Generator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetSvc()
	svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, svcFilename+".go")
	text, err := pathx.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	serviceName := strings.ToLower(stringx.From(ctx.GetServiceName().Source()).ToCamel())
	if i := strings.Index(serviceName, "service"); i > 0 {
		serviceName = strings.TrimSuffix(serviceName[:i], "-")
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"imports":     fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
		"serviceName": serviceName,
	}, fileName, false)
}
