package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sliveryou/goctl/rpc/generator"
	"github.com/sliveryou/goctl/util"
	"github.com/sliveryou/goctl/util/pathx"
)

// Client generates grpc code directly by protoc and generates
// client code by goctl.
func Client(_ *cobra.Command, args []string) error {
	protocArgs := wrapProtocCmd("protoc", args)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	source := args[0]
	grpcOutList := VarStringSliceGoGRPCOut
	goOutList := VarStringSliceGoOut
	zrpcOut := VarStringZRPCOut
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	verbose := VarBoolVerbose
	if len(grpcOutList) == 0 {
		return errInvalidGrpcOutput
	}
	if len(goOutList) == 0 {
		return errInvalidGoOutput
	}
	goOut := goOutList[len(goOutList)-1]
	grpcOut := grpcOutList[len(grpcOutList)-1]
	if len(goOut) == 0 {
		return errInvalidGrpcOutput
	}
	if len(zrpcOut) == 0 {
		return errInvalidZrpcOutput
	}
	goOutAbs, err := filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOutAbs, err := filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(goOutAbs)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(grpcOutAbs)
	if err != nil {
		return err
	}
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if !filepath.IsAbs(zrpcOut) {
		zrpcOut = filepath.Join(pwd, zrpcOut)
	}

	isGooglePlugin := len(grpcOut) > 0
	goOut, err = filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOut, err = filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	zrpcOut, err = filepath.Abs(zrpcOut)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Multiple = VarBoolMultiple
	ctx.Src = source
	ctx.GoOutput = goOut
	ctx.GrpcOutput = grpcOut
	ctx.IsGooglePlugin = isGooglePlugin
	ctx.Output = zrpcOut
	ctx.ProtocCmd = strings.Join(protocArgs, " ")
	ctx.IsGenClient = VarBoolClient
	g := generator.NewGenerator(style, verbose)
	return g.GenerateClient(&ctx)
}
