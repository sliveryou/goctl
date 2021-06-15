package upgrade

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/sliveryou/goctl/rpc/execx"
)

// Upgrade gets the latest goctl by
// go get -u github.com/sliveryou/goctl
func Upgrade(_ *cli.Context) error {
	info, err := execx.Run("GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get -u github.com/sliveryou/goctl", "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
