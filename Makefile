# 辅助工具安装列表
# 执行 go install github.com/cloudwego/hertz/cmd/hz@latest
# 执行 go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
# 执行 go install golang.org/x/tools/cmd/goimports@latest
# 执行 go install golang.org/x/vuln/cmd/govulncheck@latest
# 执行 go install mvdan.cc/gofumpt@latest
# 访问 https://golangci-lint.run/welcome/install/ 以查看安装 golangci-lint 的方法



# 默认输出帮助信息
.DEFAULT_GOAL := help
# 远程仓库
REMOTE_REPOSITORY = github.com/2451965602/LMS
# 项目 MODULE 名
MODULE = github.com/2451965602/LMS
# 当前架构
ARCH := $(shell uname -m)
PREFIX = "[Makefile]"
# 目录相关
DIR = $(shell pwd)
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output

.PHONY: help
help:
	@echo "Available targets:"
	@echo "                      Available service list: [${SERVICES}]"
	@echo "  env-up            : Start the docker-compose environment."
	@echo "  env-down          : Stop the docker-compose environment."
	@echo "  hz     : Generate Hertz scaffold based on the API IDL."
	@echo "  clean             : Remove the 'output' directories and related binaries."
	@echo "  clean-all         : Stop docker-compose services if running and remove 'output' directories and docker data."
	@echo "  fmt               : Format the codebase using gofumpt."
	@echo "  import            : Optimize import order and structure."
	@echo "  vet               : Check for possible errors with go vet."
	@echo "  lint              : Run golangci-lint on the codebase."
	@echo "  verify            : Format, optimize imports, and run linters and vet on the codebase."
	@echo "  license           : Check and add license to go file and shell script."

## --------------------------------------
## 构建与调试
## --------------------------------------

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# 生成基于 Hertz 的脚手架
.PHONY: hz-%
hz-%:
	hz update -idl ${IDL_PATH}/$*.thrift

## --------------------------------------
## 清理与校验
## --------------------------------------

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker-compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker-compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data

# 格式化代码，我们使用 gofumpt，是 fmt 的严格超集
.PHONY: fmt
fmt:
	gofumpt -l -w .

# 优化 import 顺序结构
.PHONY: import
import:
	goimports -w -local github.com/west2-online .

# 检查可能的错误
.PHONY: vet
vet:
	go vet ./...

# 代码格式校验
.PHONY: lint
lint:
	golangci-lint run --config=./.golangci.yml --tests --allow-parallel-runners --sort-results --show-stats --print-resources-usage

# 检查依赖漏洞
.PHONY: vulncheck
vulncheck:
	govulncheck ./...

.PHONY: tidy
tidy:
	go mod tidy

# 一键修正规范并执行代码检查，同时运行 license 检查
.PHONY: verify
verify: license vet fmt import lint vulncheck tidy

# 补齐 license
.PHONY: license
license:
	sh ./hack/add-license.sh

# 手动暴露环境变量
#export DOMTOK_ENVIRONMENT_STARTED=true
#export ETCD_ADDR=127.0.0.1:2379