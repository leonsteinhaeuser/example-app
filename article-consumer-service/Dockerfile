FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/article-consumer-service article-consumer-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/article-consumer-service /go/bin/article-consumer-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/article-consumer-service && \
    chmod +x /go/bin/article-consumer-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/article-consumer-service"]