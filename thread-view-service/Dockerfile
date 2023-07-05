FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/thread-view-service thread-view-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/thread-view-service /go/bin/thread-view-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/thread-view-service && \
    chmod +x /go/bin/thread-view-service

USER appuser:appuser

RUN chown appuser:appuser /go/bin/thread-view-service && chmod +x /go/bin/thread-view-service

ENTRYPOINT ["/go/bin/thread-view-service"]