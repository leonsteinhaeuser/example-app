FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/keyword-service keyword-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/keyword-service /go/bin/keyword-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/keyword-service && \
    chmod +x /go/bin/keyword-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/keyword-service"]