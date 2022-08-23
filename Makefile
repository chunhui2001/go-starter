
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

e 	?=production
c 	?=10000

tidy:
	go mod tidy

# make run e=development 
run:
	go get && GIN_ENV=$(e) go run .

# make dev
#dev:
#	go get && gin -i --appPort 8080 --port 3000 run main.go
dev:
	gin -i --appPort 8080 --port 3000 run main.go

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

# make install mod=github.com/codegangsta/gin
install:
	go install $(mod)

# show install utils
list:
	ls -alh `go env GOPATH`/bin

ngrok:
	ngrok start --config ./ngrok.yml go-starter

# make load n=10000 p=info
load:
	h2load -n$(n) -c100 -m10 --h1 "http://localhost:4000/$(p)"


