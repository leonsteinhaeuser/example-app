FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/article-service article-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/article-service /go/bin/article-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/article-service && \
    chmod +x /go/bin/article-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/article-service"]