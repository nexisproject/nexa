.PHONY: clean all build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

check_version = \
	$(if $(VERSION),,$(error 请通过 VERSION=xxx 指定版本号))

clean:
	@echo "正在清理构建文件..."
	rm -rf bin/

build-linux-amd64:
	$(call check_version)
	@echo "构建 linux/amd64 (版本: $(VERSION))"
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -trimpath -tags=sonic,poll_opt -gcflags "all=-N -l" -ldflags "-X main.Version=$(VERSION)" -o bin/nexa-linux-amd64 cmd/nexa/main.go

build-linux-arm64:
	$(call check_version)
	@echo "构建 linux/arm64 (版本: $(VERSION))"
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -trimpath -tags=sonic,poll_opt -gcflags "all=-N -l" -ldflags "-X main.Version=$(VERSION)" -o bin/nexa-linux-arm64 cmd/nexa/main.go

build-darwin-amd64:
	$(call check_version)
	@echo "构建 darwin/amd64 (版本: $(VERSION))"
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -trimpath -tags=sonic,poll_opt -gcflags "all=-N -l" -ldflags "-X main.Version=$(VERSION)" -o bin/nexa-darwin-amd64 cmd/nexa/main.go

build-darwin-arm64:
	$(call check_version)
	@echo "构建 darwin/arm64 (版本: $(VERSION))"
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -trimpath -tags=sonic,poll_opt -gcflags "all=-N -l" -ldflags "-X main.Version=$(VERSION)" -o bin/nexa-darwin-arm64 cmd/nexa/main.go

build-windows-amd64:
	$(call check_version)
	@echo "构建 windows/amd64 (版本: $(VERSION))"
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -trimpath -tags=sonic,poll_opt -gcflags "all=-N -l" -ldflags "-X main.Version=$(VERSION)" -o bin/nexa-windows-amd64.exe cmd/nexa/main.go

all: clean build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
	@echo "全部平台构建完成"
