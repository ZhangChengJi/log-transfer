FROM registry.cn-qingdao.aliyuncs.com/nqkj-snapshot/golang:1.15 as build

MAINTAINER 380702562@qq.com


# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN  go mod download

WORKDIR /go/release
ADD . .
# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o  log-transfer main.go

FROM scratch as prod

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build /go/release/log-transfer /

CMD ["./log-transfer"]