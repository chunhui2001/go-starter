
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

### websocket 
# https://hoohoo.top/blog/20220320172715-go-websocket/

$ go get github.com/gorilla/websocket

### html template
# https://betterprogramming.pub/how-to-render-html-pages-with-gin-for-golang-9cb9c8d7e7b6

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

