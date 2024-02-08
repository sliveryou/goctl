package ctx

import (
	"bytes"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sliveryou/goctl/rpc/execx"
	"github.com/sliveryou/goctl/util/pathx"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestProjectFromGoMod(t *testing.T) {
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

	ctx, err := projectFromGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, projectName, ctx.Path)
	assert.Equal(t, dir, ctx.Dir)
}

func Test_getRealModule(t *testing.T) {
	type args struct {
		workDir string
		execRun execx.RunFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *Module
		wantErr bool
	}{
		{
			name: "single module",
			args: args{
				workDir: "/home/foo",
				execRun: func(arg, dir string, in ...*bytes.Buffer) (string, error) {
					return `{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.16"}`, nil
				},
			},
			want: &Module{
				Path:      "foo",
				Dir:       "/home/foo",
				GoMod:     "/home/foo/go.mod",
				GoVersion: "go1.16",
			},
		},
		{
			name: "go work multiple modules",
			args: args{
				workDir: "/home/bar",
				execRun: func(arg, dir string, in ...*bytes.Buffer) (string, error) {
					return `
					{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.19"
					}
					{
						"Path":"bar",
						"Dir":"/home/bar",
						"GoMod":"/home/bar/go.mod",
						"GoVersion":"go1.19"
					}`, nil
				},
			},
			want: &Module{
				Path:      "bar",
				Dir:       "/home/bar",
				GoMod:     "/home/bar/go.mod",
				GoVersion: "go1.19",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRealModule(tt.args.workDir, tt.args.execRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRealModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRealModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodePackages(t *testing.T) {
	tests := []struct {
		name    string
		data    io.Reader
		want    []Module
		wantErr bool
	}{
		{
			name: "single module",
			data: strings.NewReader(`{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.16"}`),
			want: []Module{
				{
					Path:      "foo",
					Dir:       "/home/foo",
					GoMod:     "/home/foo/go.mod",
					GoVersion: "go1.16",
				},
			},
		},
		{
			name: "go work multiple modules",
			data: strings.NewReader(`
					{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.19"
					}
					{
						"Path":"bar",
						"Dir":"/home/bar",
						"GoMod":"/home/bar/go.mod",
						"GoVersion":"go1.19"
					}`),
			want: []Module{
				{
					Path:      "foo",
					Dir:       "/home/foo",
					GoMod:     "/home/foo/go.mod",
					GoVersion: "go1.19",
				},
				{
					Path:      "bar",
					Dir:       "/home/bar",
					GoMod:     "/home/bar/go.mod",
					GoVersion: "go1.19",
				},
			},
		},
		{
			name: "There are extra characters at the beginning",
			data: strings.NewReader(`Active code page: 65001
					{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.19"
					}`),
			want: []Module{
				{
					Path:      "foo",
					Dir:       "/home/foo",
					GoMod:     "/home/foo/go.mod",
					GoVersion: "go1.19",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodePackages(tt.data)
			if err != nil {
				t.Errorf("decodePackages() error %v,wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("decodePackages() = %v,want  %v", result, tt.want)

			}
		})
	}
}
