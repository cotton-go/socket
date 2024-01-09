#socket

### 目录结构
```
.
├── config
│   └── local.yml
├── go.mod
├── go.sum
├── init.go
├── pkg
│   ├── app
│   │   ├── app.go
│   │   └── option.go
│   ├── cache
│   │   ├── cache.go
│   │   ├── memory.go
│   │   └── redis.go
│   ├── client
│   │   └── client.go
│   ├── codec
│   │   ├── AES_CBC.go
│   │   ├── AES_ECB.go
│   │   ├── codec.go
│   │   ├── default.go
│   │   ├── DES_CBC.go
│   │   ├── DES_ECB.go
│   │   └── event.go
│   ├── config
│   │   └── config.go
│   ├── connection
│   │   ├── gonnection.go
│   │   └── option.go
│   ├── event
│   │   ├── constant.go
│   │   └── event.go
│   ├── log
│   │   └── log.go
│   ├── registry
│   │   └── registry.go
│   ├── server
│   │   ├── grpc
│   │   │   ├── config.go
│   │   │   ├── grpc.go
│   │   │   └── option.go
│   │   ├── http
│   │   │   ├── config.go
│   │   │   └── http.go
│   │   ├── server.go
│   │   └── tcp
│   │       ├── config.go
│   │       ├── option.go
│   │       ├── tcp.go
│   │       └── tcp_test.go
│   ├── snowflake
│   │   ├── init.go
│   │   └── snowflake.go
│   └── worker
│       ├── option.go
│       ├── worker.go
│       └── worker_test.go
├── README.md
├── server.go
└── tests
    ├── client
    │   └── client.go
    └── server
        └── server.go
```