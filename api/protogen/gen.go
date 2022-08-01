package protogen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/urfave/cli"

	"github.com/sliveryou/goctl/api/parser"
	"github.com/sliveryou/goctl/util"
)

// ProtoCommand gen proto file from command line
func ProtoCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")

	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	return DoGenProto(apiFile, dir)
}

// DoGenProto gen proto file with api file
func DoGenProto(apiFile, dir string) error {
	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	apiBase := filepath.Base(apiFile)
	apiName := apiBase[:len(apiBase)-len(filepath.Ext(apiBase))]

	logx.Must(util.MkdirIfNotExist(dir))
	f, err := os.Create(path.Join(dir, apiName+"-rpc.proto"))
	logx.Must(err)
	defer f.Close()

	ts, err := BuildTypes(api)
	logx.Must(err)

	rs, hasEmpty := BuildRPCs(api)

	_, err = f.WriteString("syntax = \"proto3\";\n\npackage rpc;")
	logx.Must(err)

	if hasEmpty {
		_, err = f.WriteString("\n\n// Empty 空消息\nmessage Empty {\n}")
		logx.Must(err)
	}

	_, err = f.WriteString("\n\n" + ts + "\n\n" + rs)
	logx.Must(err)

	fmt.Println(aurora.Green("Done."))

	return nil
}
