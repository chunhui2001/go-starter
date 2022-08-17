
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

e 	?=production
c 	?=10000

# make dev
dev:
	go get && gin -i --appPort 8080 --port 3000 run main.go

# make run e=development 
run:
	go get && GIN_ENV=$(e) go run .

build:
	docker rmi -f go-starter:1.0 && docker build . -t go-starter:1.0

up:
	docker-compose -f docker-compose.yml up -d

log:
	docker logs -f --tail 1000 go-starter

rm:
	docker rm -f go-starter

# make load n=10000 p=info
load:
	h2load -n$(n) -c100 -m10 --h1 "http://localhost:4000/$(p)"

clear:
	go clean --modcache
