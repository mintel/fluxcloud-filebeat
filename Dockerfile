FROM golang:1.13-alpine as base

RUN apk add --no-cache make git g++ musl-dev linux-headers bash ca-certificates
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base AS builder
COPY . .
RUN GOOS=linux GOARCH=amd64 make

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/fluxcloud-filebeat /usr/local/bin/fluxcloud-filebeat
COPY --from=builder /app/templates /app/templates

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/fluxcloud-filebeat"]
