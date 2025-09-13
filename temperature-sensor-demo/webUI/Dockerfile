FROM golang:1.22.9-alpine3.19 AS builder

WORKDIR /build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY . .

RUN go get github.com/kubeedge/beehive@latest

RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags="-s -w" -o main ./main.go

FROM ubuntu:22.04

WORKDIR /webUI

COPY --from=builder /build/main ./main

COPY ./static ./static

EXPOSE 8080

CMD ["./main"]