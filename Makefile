version := $(shell /bin/date "+%Y-%m-%d %H:%M")

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" goctl.go
	$(if $(shell command -v upx), upx goctl)
mac:
	GOOS=darwin go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl-darwin goctl.go
	$(if $(shell command -v upx), upx goctl-darwin)
win:
	GOOS=windows go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl.exe goctl.go
	$(if $(shell command -v upx), upx goctl.exe)
linux:
	GOOS=linux go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o goctl-linux goctl.go
	$(if $(shell command -v upx), upx goctl-linux)
fmt:
	@find . -name '*.go' -not -path "./vendor/*" -not -name "*.pb.go" | xargs gofumpt -w -s -extra
	@find . -name '*.go' -not -path "./vendor/*" -not -name "*.pb.go" | xargs -n 1 -I {} -t goimports-reviser -file-path {} -local "github.com/tal-tech" project-name "github.com/sliveryou/goctl" -rm-unused
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shfmt -w -s -i 2 -ci -bn -sr
