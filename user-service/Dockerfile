FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/user-service user-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/user-service /go/bin/user-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/user-service && \
    chmod +x /go/bin/user-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/user-service"]