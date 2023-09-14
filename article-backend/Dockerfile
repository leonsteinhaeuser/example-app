FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/article-backend article-backend/main.go

FROM alpine:latest

COPY --from=builder /go/bin/article-backend /go/bin/article-backend

RUN adduser appuser -u 1234 --disabled-password && \
    chown -R appuser:appuser /go/bin/article-backend && \
    chmod +x /go/bin/article-backend

USER appuser:appuser

ENTRYPOINT ["/go/bin/article-backend"]
