
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

e 		?=local
c 		?=10000
#zone 	?=Asia/Shanghai
zone 	?=UTC
WSS_HOST	?=ws://127.0.0.1:8080
APP_PORT 	?=8080
GIT_HASH 	?=$(shell git rev-parse HEAD)
COMMITER 	?=$(shell git log --format="%H" -n 1 | grep Author| cut -d ' ' -f3- | sed 's/[\<\>]*//g')
TIME 		?=$(shell date +%s)
#GOOS 		?=darwin
#GOOS 		?=windows
GOOS 		?=linux

### 整理模块
# 确保go.mod与模块中的源代码一致。
# 它添加构建当前模块的包和依赖所必须的任何缺少的模块，删除不提供任何有价值的包的未使用的模块。
# 它也会添加任何缺少的条目至go.mod并删除任何不需要的条目。
# make tidy
tidy:
	go mod tidy

### 安装模块
# make install mod=github.com/codegangsta/gin
install:
	go get github.com/codegangsta/gin
	go install github.com/codegangsta/gin
	go get $(mod)
	go install $(mod)

### 下载模块
get:
	go get

### 启动开发程序
# make run e=development 
run:
	rm -rf gin-bin
	GIN_ENV=$(e) go run .

### 启动调试程序, 当代码变化时自动重启
# make dev
dev:
	TZ=$(zone) GIN_ENV=$(e) WSS_HOST=$(WSS_HOST) gin -i --appPort 8080 --port 3000 run main.go

### 构建程序镜像
build: Built
	docker rmi -f go-starter:1.0 && docker build . -t go-starter:1.0  -m 4g

### 通过容器启动
up: rm
	docker-compose -f docker-compose.yml up -d

### 查看程序日志
logs:
	docker logs -f --tail 1000 go-starter

### 删除程序容器
rm:
	docker rm -f go-starter
	rm -rf ./dist

### 构建跨平台的可执行程序
Built:
	env GOOS=$(GOOS) GOARCH=amd64 go build -buildvcs -ldflags "-X main.Name=go-starter -X main.Author=$(COMMITER) -X main.Commit=$(GIT_HASH) -X main.Time=$(TIME)" -o ./dist/go-starter-native-$(GOOS)-amd64 ./main.go
	
### 删除所有缓存的依赖包
# clear modcache
clear:
	go clean --modcache
	rm -rf `go env GOPATH`/bin/go-starter
	@#rm -rf `go env GOPATH`/bin/*
	rm -rf dist gin-bin

### 显示已安装的可执行程序
# show install utils
list:
	ls -alh `go env GOPATH`/bin

# make ngrok
#ngrok:
#	ngrok start --config ./ngrok.yml go-starter

### 性能测试
# make load n=10000 p=info
load:
	h2load -n$(n) -c100 -m10 --h1 "http://localhost:4000/$(p)"


