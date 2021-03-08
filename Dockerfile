FROM golang:1.15-alpine3.13 as builder

MAINTAINER 380702562@qq.com
RUN adduser -u 10001 -D app-runner
# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod .
COPY go.sum .
RUN go mod download

WORKDIR /app


# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o log-transfer .

FROM scratch
WORKDIR /app
COPY --from=builder /build/log-transfer  /app/



EXPOSE 9000
USER app-runner
ENTRYPOINT ["./log-transfer"]