FROM golang:1.23.3-alpine3.19 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags="-s -w" -o main ./main.go

FROM alpine:3.19

COPY --from=builder /app/main .

EXPOSE 5502

CMD ["./main"]

