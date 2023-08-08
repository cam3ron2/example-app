FROM golang:latest as builder
COPY src/ /go/src/
WORKDIR /go/src/
RUN export ARCH=$(case $(uname -m) in aarch64|arm64) echo arm64;; x86_64) echo amd64;; esac); \
  export OS=$(echo $(uname -s) | tr '[:upper:]' '[:lower:]'); \
  env GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=0 go build -o /go/example-app .

FROM alpine:3.18.3
COPY --from=builder /go/example-app /usr/bin/
RUN chmod +x /usr/bin/example-app
WORKDIR /app
USER nobody
ENTRYPOINT [ "/usr/bin/example-app" ]