
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

e 		?=local
c 		?=10000
#zone 	?=Asia/Shanghai
zone 	?=UTC
WSS_HOST	?=ws://127.0.0.1:8080

# make tidy
tidy:
	go mod tidy

# make install mod=github.com/codegangsta/gin
install:
	go install $(mod)

# make get
get:
	go get

# make run e=development 
run:
	GIN_ENV=$(e) go run .

# make dev
dev:
	TZ=$(zone) GIN_ENV=$(e) WSS_HOST=$(WSS_HOST) gin -i --appPort 8080 --port 3000 run main.go

# build docker image
build:
	docker rmi -f go-starter:1.0 && docker build . -t go-starter:1.0

# docker up
up: rm
	docker-compose -f docker-compose.yml up -d

# view logs
logs:
	docker logs -f --tail 1000 go-starter

# delete container
rm:
	docker rm -f go-starter

# clear modcache
clear:
	go clean --modcache
	rm -rf `go env GOPATH`/bin/go-starter

# show install utils
list:
	ls -alh `go env GOPATH`/bin

# make ngrok
ngrok:
	ngrok start --config ./ngrok.yml go-starter

# make load n=10000 p=info
load:
	h2load -n$(n) -c100 -m10 --h1 "http://localhost:4000/$(p)"


