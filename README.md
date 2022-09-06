# Example App

![codeql](https://github.com/cam3ron2/example-app/actions/workflows/codeql.yml/badge.svg)
[![tests](https://github.com/cam3ron2/example-app/actions/workflows/tests.yml/badge.svg)](https://github.com/cam3ron2/example-app/actions/workflows/tests.yml)
<!-- [![build](https://github.com/cam3ron2/example-app/actions/workflows/build-app.yml/badge.svg)](https://github.com/cam3ron2/example-app/actions/workflows/build-app.yml) -->
<!-- ![downloads](https://img.shields.io/github/downloads/cam3ron2/example-app/v1.0.6/total) -->
![License](https://img.shields.io/github/license/cam3ron2/example-app)
![Language](https://img.shields.io/badge/language-Go-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/cam3ron2/example-app/src)](https://goreportcard.com/report/github.com/cam3ron2/example-app/src)

A lightweight golang webserver and client useful for testing and demos.

Starting the webserver:

```bash
$ exmaple-app server --help
Starts a server instance

Usage:
  example-app server [flags]

Flags:
  -d, --delay int         response delay in ms
  -f, --fail int          % of requests to fail, ex 10 = 10%
  -F, --health-fail int   % of requests to /healthz to fail, ex 10 = 10%
  -h, --help              help for server
  -p, --port int          port to listen on (default 8080)

Global Flags:
  -D, --datadog   Enable DataDog trace collection

$ example-app server
[Server] 2022/08/10 13:00:34 Starting Server on port :8080
[Server] 2022/08/10 13:00:34 Server is ready to handle requests at :8080
```

Starting the client worker:

```bash
$ example-app worker --help
Starts a worker instance

Usage:
  example-app worker [flags]

Flags:
  -f, --fail int          % of requests to fail, ex 10 = 10%
  -F, --health-fail int   % of requests to /healthz to fail, ex 10 = 10%
  -P, --health-port int   worker healthcheck Port (default 8081)
  -h, --help              help for worker
  -p, --port int          target port (default 8080)
  -r, --rate int          rate of requests per second (default 1)
  -u, --url string        target URL (default "http://localhost")

Global Flags:
  -D, --datadog   Enable DataDog trace collection

$ example-app worker
[Worker] 2022/08/10 13:02:31 Starting Worker on port :8081
[Worker] 2022/08/10 13:02:31 Server is ready to handle requests at :8081
```

## DataDog Configuration

## TODO

- Add completion
- Finish writing tests

