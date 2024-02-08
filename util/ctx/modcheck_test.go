package ctx

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/sliveryou/goctl/rpc/execx"
	"github.com/sliveryou/goctl/util/pathx"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestIsGoMod(t *testing.T) {
	// create mod project
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := pathx.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	_, err = execx.Run("go mod init "+projectName, dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	isGoMod, err := IsGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, true, isGoMod)
}

func TestIsGoModNot(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := pathx.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	isGoMod, err := IsGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, false, isGoMod)
}

func TestIsGoModWorkDirIsNil(t *testing.T) {
	_, err := IsGoMod("")
	assert.Equal(t, err.Error(), func() string {
		return "the work directory is not found"
	}())
}
