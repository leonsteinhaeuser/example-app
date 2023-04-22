FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/view-service view-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/view-service /go/bin/view-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/view-service && \
    chmod +x /go/bin/view-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/view-service"]