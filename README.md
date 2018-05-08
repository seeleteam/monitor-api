# Monitor API

## Env
>recommend
```text
go 1.10+
```

>To use config file `app.conf`, you can choose the following ways
- export MONITOR_API_CONFIG_PATH=<MONITOR_API_CONFIG_PATH>
- app.conf
- conf/app.conf

## Project structure
```text
├── api: api interface
│   ├── filters:  request filter
│   ├── handlers: router handler
│   └── routers:  the http router
├── conf: config files
├── config: config struct
├── core
│   ├── config: configure
│   ├── logs: third logger
│   └── utils: utils
├── rpc: json rpc
├── server: monitor server
├── tools:  build script
│   └── prebuild
│       └── shell
├── vendor: third dependencies
└── ws: web socket

```

## Start
```bash
cd monitor-api
# generate the executable file
bash build.sh buildd
# execute the file
bin/monitor-api
```

## Config
```text
run_mode = dev

# http server address, format ip:port
addr = :9997

LimitConnection = 0
# whether or not to use grace shutdown 
Graceful = false
# graceful shutdown time, unit seconds
#DefaultHammerTime = 30

# enable web socket
EnableWebSocket = true

# enable rpc
EnableRPC = true
DisableConsoleColor = false

# enable write log out
WriteLog = true

# websocket server address, format ip:port
WsURL = :9997

# web socket api
WsRouter = /api

# every 10s send the node info to monitor server
WsFullEventTickerTime = 10
# every 2s send the block info, if the block height changed
WsLatestBlockEventTickerTime = 2

# if web socket occur error, reconnect delay 5s
DelayReConnTime = 5

# if rpc occur error, reconnetct and resend delay 5s
DelaySendTime = 5

# if rpc server occur error over 10, report error to monitor server
ReportErrorAfterTimes = 10

# RPC server addr for go-seele node, format ip:port
RPCURL = 127.0.0.1:55027

# log level, debug, info, warn, error, fatal, panic
LogLevel = debug

```