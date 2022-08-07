
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

e 	?=production
c 	?=10000

# make run e=development 
run:
	go get && GIN_ENV=$(e) go run .

build:
	docker build . -t go-starter:1.0

up:
	docker-compose -f docker-compose.yml up -d

log:
	docker logs -f go-starter

rm:
	docker rm -f go-starter

# make load c=10000 
load:
	h2load -n$(c) -c100 -m10 --h1 "http://localhost:4000/info"
