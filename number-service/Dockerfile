FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/number-service number-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/number-service /go/bin/number-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/number-service && \
    chmod +x /go/bin/number-service

USER appuser:appuser

RUN chown appuser:appuser /go/bin/number-service && chmod +x /go/bin/number-service

ENTRYPOINT ["/go/bin/number-service"]