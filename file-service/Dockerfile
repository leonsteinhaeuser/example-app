FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/file-service file-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/file-service /go/bin/file-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/file-service && \
    chmod +x /go/bin/file-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/file-service"]