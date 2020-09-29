FROM golang:1.12-alpine3.10
RUN apk add gcc libc-dev

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io/

WORKDIR $GOPATH/src/mapper

COPY . .
RUN go build .
ENTRYPOINT ["./main.go"]






