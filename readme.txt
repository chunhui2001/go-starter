### 限制
Content-Type restriction

### 版本列表
https://go.dev/doc/devel/release

## All releases
https://go.dev/dl/

### 设置国内代理
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct

### 查看文件size
stat -f%z README.md 

### Install
$ go get github.com/codegangsta/gin

### Refer
# https://go.dev/doc/tutorial/web-service-gin

### set gopath
> 21 export GOROOT=/usr/local/go
> 22 export GOPATH=$HOME/go
> 27 export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

### steps:
$ mkdir go-starter && cd go-starter
$ go mod init go-starter

$ go get -u github.com/codegangsta/gin
$ go get -u github.com/jteeuwen/go-bindata/...
$ go get -u github.com/elazarl/go-bindata-assetfs/...

$ go get && go run .

### 查看当前进程打开的文件描述符
$ lsof -p 1243 | wc -l

### 以线程模式查看下进程 31951 的所有线程情况
ps -T 31951 

### 查看系统设置
$ ulimit -a 

### 
$ ulimit -Sn

### 设置
$ ulimit -n 1000000
OR
$ sysctl -w fs.file-max=1000000

### and /etc/security/limits.conf or /etc/sysctl.conf change:
fs.file-max = 1000000


# binary will be $(go env GOPATH)/bin/golangci-lint
$ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
$ golangci-lint --version

### Go Playground
https://go.dev/play/

### building microservices go gin
https://blog.logrocket.com/building-microservices-go-gin/

### websocket 
# https://hoohoo.top/blog/20220320172715-go-websocket/

### A Million WebSockets and Go
https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/

### How to use websockets in Golang: best tools and step-by-step guide
https://yalantis.com/blog/how-to-build-websockets-in-go/

### Handle 'connection reset by peer' error in Go
https://gosamples.dev/connection-reset-by-peer/

$ go get github.com/gorilla/websocket

### html template
# https://betterprogramming.pub/how-to-render-html-pages-with-gin-for-golang-9cb9c8d7e7b6
# http://2016.8-p.info/post/06-18-go-html-template/

### goview
# https://github.com/foolin/goview
# https://curatedgo.com/r/goview-is-a-foolingoview/index.html
# https://www.godoc.org/github.com/foolin/goview

### Gin doc
https://chenyitian.gitbooks.io/gin-web-framework/content/docs/14.html

### utils
https://github.com/ubiq/go-ubiq
https://pkg.go.dev/github.com/ubiq/go-ubiq@v3.0.1+incompatible

### 31个！Golang常用工具
https://blog.csdn.net/QcloudCommunity/article/details/126047057

### Session Cookie Authentication in Golang (With Complete Examples)
https://www.sohamkamani.com/golang/session-cookie-authentication/

### 我用 Golang 的 Gin/bindata (+React) 尝试了一个二进制文件
https://qiita.com/wadahiro/items/4173788d54f028936723

### Build CRUD RESTful API Server with Golang, Gin, and MongoDB
https://codevoweb.com/crud-restful-api-server-with-golang-and-mongodb/

### gin-rate-limit
https://github.com/JGLTechnologies/gin-rate-limit

### How To Make HTTP Requests in Go
https://www.digitalocean.com/community/tutorials/how-to-make-http-requests-in-go

### http2curl
https://github.com/moul/http2curl

### log to kafka
https://github.com/gfremex/logrus-kafka-hook

### Gin binding in Go: A tutorial with examples
https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

### Embedding Git Commit Information in Go Binaries
https://icinga.com/blog/2022/05/25/embedding-git-commit-information-in-go-binaries/

### https://github.com/golang-standards/project-layout
https://github.com/golang-standards/project-layout

### Ratago is a (mostly-compliant) implementation of an XSLT 1.0 processor written in Go and released under an MIT license.
https://github.com/jbowtie/ratago
https://github.com/wamuir/go-xslt

### sqlx is a library which provides a set of extensions on go's standard database/sql library. 
https://github.com/jmoiron/sqlx

### How to Implement MySQL Transactions with Golang On Linux Server
vultr.com/ja/docs/how-to-implement-mysql-transactions-with-golang-on-linux-server/

### Setting a time limit for a transaction in MySQL/InnoDB
https://serverfault.com/questions/241823/setting-a-time-limit-for-a-transaction-in-mysql-innodb

### Executing transactions
https://go.dev/doc/database/execute-transactions

### Pointers
// https://www.golang-book.com/books/intro/8

### golang-tls
https://gist.github.com/denji/12b3a568f092ab951456

### xslt
https://www.ardanlabs.com/blog/2013/11/using-xslt-with-go.html

### excel
https://github.com/tealeg/csv2xlsx/blob/master/main.go

### Installing protoc
http://google.github.io/proto-lens/installing-protoc.html

### How To Bulk Index Elasticsearch Documents Using Golang
https://kb.objectrocket.com/elasticsearch/how-to-bulk-index-elasticsearch-documents-using-golang-450
https://onexlab-io.medium.com/elasticsearch-bulk-insert-json-data-322f97115d8d

>>>> $ PROTOC_ZIP=protoc-3.14.0-osx-x86_64.zip
>>>> $ curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
>>>> $ sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
>>>> $ sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
>>>> $ rm -f $PROTOC_ZIP

$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

### Go Protocol Buffer Tutorial
https://tutorialedge.net/golang/go-protocol-buffer-tutorial/

### Go package to make lightweight ASCII line graph in command line apps with no other dependencies
https://github.com/guptarohit/asciigraph
https://golangexample.com/go-package-to-make-lightweight-ascii-line-graph-in-command-line-apps-with-no-other-dependencies/
> ping -i.2 google.com | grep -oP '(?<=time=).*(?=ms)' --line-buffered | /asciigraph -r -h 10 -w 40 -c "realtime plot data (google ping in ms) from stdin"

### Profile-specific application properties (application-{profile}.properties and YAML variants).
### Properties are considered in the following order:
> 1. Command line arguments.
> 2. OS environment variables.
> 3. RandomValuePropertySource that has properties only in random.*.
> 4. Profile-specific application properties (application-{profile}.properties and YAML variants).
> 5. Application properties (application.properties and YAML variants).
> 6. Default properties
https://github.com/furkilic/go-boot-config

### 通过实例理解Go标准库http包是如何处理keep-alive连接的
https://tonybai.com/2021/01/08/understand-how-http-package-deal-with-keep-alive-connection/

### GraphQL api server using golang Gin framework
https://www.agiliq.com/blog/2020/04/graphql-api-server-using-golang-gin-framework/

### gqlgen
https://gqlgen.com/config/

### Mapping GraphQL scalar types to Go types
https://gqlgen.com/reference/scalars/

### graphql tools
https://hygraph.com/blog/graphql-tools

### A helper to merge structs and maps in Golang. Useful for configuration default values, avoiding messy if-statements.
https://github.com/imdario/mergo

### A Complete Guide to JSON in Golang (With Examples)
https://www.sohamkamani.com/golang/json/

### GoLang: When to use string pointers
https://dhdersch.github.io/golang/2016/01/23/golang-when-to-use-string-pointers.html

### Go Fundamentals: Arrays and Slices (+caveats of appending)
https://www.integralist.co.uk/posts/go-slices/

### Arrays, slices (and strings): The mechanics of 'append'
https://go.dev/blog/slices

### HTML Versus XHTML
https://www.w3schools.com/html/html_xhtml.asp

### Gin 101: Enable CSRF middleware
https://medium.com/@pointgoal/gin-101-enable-csrf-middleware-27faa2420186

## Google Doc
https://docs.google.com/

## Google API Console
# https://console.developers.google.com/

## 快速开始
# https://developers.google.com/sheets/api/quickstart/go

## 修改文件权限
# https://developers.google.com/drive/api/v3/reference/revisions

### (Go) Generate an AWS (S3) Pre-Signed URL using Signature V4
# https://www.example-code.com/golang/aws_pre_signed_url_v4.asp

### How to download only the beginning of a large file with go http.Client (and app engine urlfetch)
https://stackoverflow.com/questions/27844307/how-to-download-only-the-beginning-of-a-large-file-with-go-http-client-and-app

### Recombining large chunked zip download in GO
https://stackoverflow.com/questions/53605400/recombining-large-chunked-zip-download-in-go

### Downloading large files in Go
http://cavaliercoder.com/blog/downloading-large-files-in-go.html

### Introducing Ristretto: A High-Performance Go Cache
https://dgraph.io/blog/post/introducing-ristretto-high-perf-go-cache/

### github.com/robfig/cron
https://pkg.go.dev/github.com/robfig/cron/v3#section-readme
https://en.wikipedia.org/wiki/Cron
https://crontab.guru/

### Golang Cron V3 Timed Tasks
https://www.sobyte.net/post/2021-06/golang-cron-v3-timed-tasks/

### OpenSearch 2.6.x
https://opensearch.org/docs/latest/install-and-configure/install-opensearch/docker/

-----------------------------------------------------------------
GOOS - Target Operating System		GOARCH - Target Platform
-----------------------------------------------------------------
android								arm
darwin								386
darwin								amd64
darwin								arm
darwin								arm64
dragonfly							amd64
freebsd								386
freebsd								amd64
freebsd								arm
linux								386
linux								amd64
linux								arm
linux								arm64
linux								ppc64
linux								ppc64le
linux								mips
linux								mipsle
linux								mips64
linux								mips64le
netbsd								386
netbsd								amd64
netbsd								arm
openbsd								386
openbsd								amd64
openbsd								arm
plan9								386
plan9								amd64
solaris								amd64
windows								386
windows								amd64

### 【golang】性能优化
https://blog.csdn.net/shanxiaoshuai/article/details/121720800

### Parallel For-Loop
http://www.golangpatterns.info/concurrency/parallel-for-loop

### Vec中的Rust元组：使用val0方法时出现编译错误
https://www.javaroad.cn/questions/341960

### Golang性能调优(go-torch, go tool pprof)
http://www.manongjc.com/detail/51-ygrewovmblzibqg.html

### fasthttp
https://github.com/valyala/fasthttp

### Mac系统上安装AB工具
https://blog.csdn.net/qq_42700121/article/details/120945354
https://www.jianshu.com/p/a7ee2ffb5c0f

-n 请求树
-c 并发数（访问人数）
-t 请求时间最大数

$ ab -n 1 -c 1 http://www.baidu.com/


