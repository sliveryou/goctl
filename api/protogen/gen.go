package protogen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/sliveryou/goctl/api/parser"
	"github.com/sliveryou/goctl/util/pathx"
)

var (
	// VarStringDir describes the directory.
	VarStringDir string
	// VarStringAPI describes the API.
	VarStringAPI string
	// VarStringRemoveBeforeDelimiter describes the delimiter.
	VarStringRemoveBeforeDelimiter string
)

// ProtoCommand gen proto file from command line
func ProtoCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir

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

	logx.Must(pathx.MkdirIfNotExist(dir))
	f, err := os.Create(path.Join(dir, apiName+".proto"))
	logx.Must(err)
	defer f.Close()

	ts, err := BuildTypes(api)
	logx.Must(err)

	rs, hasEmpty := BuildRPCs(api, apiName)

	_, err = f.WriteString(fmt.Sprintf("syntax = \"proto3\";\n\noption go_package = \"./pb\";\n\npackage %s;", apiName))
	logx.Must(err)

	_, err = f.WriteString("\n\n" + rs)
	logx.Must(err)

	if hasEmpty {
		_, err = f.WriteString("\n\n// Empty 空消息\nmessage Empty {\n}")
		logx.Must(err)
	}

	_, err = f.WriteString("\n\n" + ts)
	logx.Must(err)

	color.Green.Println("Done.")

	return nil
}
