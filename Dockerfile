ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git --no-cache g++ gcc libxslt-dev libxml2-dev && rm -rf /var/cache/apk/*

RUN mkdir -p /dist
WORKDIR /dist

### COPY . .
### RUN make Built GOOS=linux && ln -s go-starter-native-linux-amd64 app
COPY ./dist/go-starter-native-linux-amd64 ./app


FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN apk add -U tzdata && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone

RUN mkdir -p /dist
WORKDIR /dist
COPY --from=builder /dist/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
