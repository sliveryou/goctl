package gateway

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/sliveryou/goctl/internal/cobrax"
	"github.com/sliveryou/goctl/util/ctx"
	"github.com/sliveryou/goctl/util/pathx"
	"github.com/spf13/cobra"
)

var (
	varStringHome   string
	varStringRemote string
	varStringBranch string
	varStringDir    string

	Cmd = cobrax.NewCommand("gateway", cobrax.WithRunE(generateGateway))
)

func init() {
	Cmd.PersistentFlags().StringVar(&varStringHome, "home")
	Cmd.PersistentFlags().StringVar(&varStringRemote, "remote")
	Cmd.PersistentFlags().StringVar(&varStringBranch, "branch")
	Cmd.PersistentFlags().StringVar(&varStringDir, "dir")
}

func generateGateway(*cobra.Command, []string) error {
	if err := pathx.MkdirIfNotExist(varStringDir); err != nil {
		return err
	}

	if _, err := ctx.Prepare(varStringDir); err != nil {
		return err
	}

	etcContent, err := pathx.LoadTemplate(category, etcTemplateFileFile, etcTemplate)
	if err != nil {
		return err
	}

	mainContent, err := pathx.LoadTemplate(category, mainTemplateFile, mainTemplate)
	if err != nil {
		return err
	}

	etcDir := filepath.Join(varStringDir, "etc")
	if err := pathx.MkdirIfNotExist(etcDir); err != nil {
		return err
	}
	etcFile := filepath.Join(etcDir, "gateway.yaml")
	if err := os.WriteFile(etcFile, []byte(etcContent), 0644); err != nil {
		return err
	}

	mainFile := filepath.Join(varStringDir, "main.go")
	return os.WriteFile(mainFile, []byte(mainContent), 0644)
}
