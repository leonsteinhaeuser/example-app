FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/thread-service thread-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/thread-service /go/bin/thread-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/thread-service && \
    chmod +x /go/bin/thread-service

USER appuser:appuser

RUN chown appuser:appuser /go/bin/thread-service && chmod +x /go/bin/thread-service

ENTRYPOINT ["/go/bin/thread-service"]