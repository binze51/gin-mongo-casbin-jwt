VERSION := $(shell git describe --tags --always --dirty="-dev")
GIT_COMMIT := $(shell git rev-list -1 HEAD)

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

version:
	@echo '输出项目版本'
	@echo $(VERSION) \(git commit: $(GIT_COMMIT)\)

test:
	@echo '运行项目单元测试'
	go mod tidy\
	&& go test -v

binary:
	@echo '构建linux二进制文件'
	go mod tidy\
	&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binzekeji:$(VERSION)