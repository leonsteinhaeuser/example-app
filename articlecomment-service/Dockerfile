FROM golang:alpine AS builder

WORKDIR /go/src/app
ADD . .

RUN go build -o /go/bin/articlecomment-service articlecomment-service/main.go

FROM alpine:latest

COPY --from=builder /go/bin/articlecomment-service /go/bin/articlecomment-service

RUN adduser appuser --disabled-password && \
    chown -R appuser:appuser /go/bin/articlecomment-service && \
    chmod +x /go/bin/articlecomment-service

USER appuser:appuser

ENTRYPOINT ["/go/bin/articlecomment-service"]