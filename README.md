# Monitor API

## Env

>recommend

```text
go 1.10+
```

## Project structure

```text
┌── api: api interface
│   ├── filters:  request filter
│   ├── handlers: router handler
│   └── routers:  the http router
│── cmd: command module
├── config: config struct
├── core
│   ├── config: configure
│   ├── logs: third logger
│   └── utils: utils
├── rpc: json rpc
├── server: monitor server
├── vendor: third dependencies
└── ws: web socket

```

## Start

```bash
cd monitor-api
# generate the executable file
make

# in windows you should use ./buildall.bat replace

# execute the file
./monitor-api start

or

./monitor-api start -c <configfile>
```

## Warn

- default app.conf and monitor.json should be in *the same config dir*, the structure should be like

```text
┌── monitor-api
│   └── config
│       ├── app.conf
│       └── monitor.json
```

- if set env `MONITOR_CONFIG_FILE` the MonitorConfigFile in `app.conf` will be override!

## Docker

### build

```bash
docker build --rm -f Dockerfile -t monitor-api:latest .
```

### run

```bash
docker run -v <your app.conf directory>:/monitor-api/config:ro -it monitor-api
```

## Config

```text
run_mode = dev

# http server address, format ip:port
addr = :9997

# the path store temp files, like log
TempFolder = ""

LimitConnection = 0

# enable web socket
EnableWebSocket = true

# enable rpc
EnableRPC = true
DisableConsoleColor = false

# enable write log out
WriteLog = true

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