# Go builder container
FROM golang:alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go/src/github.com/seeleteam/monitor-api

RUN cd /go/src/github.com/seeleteam/monitor-api && make all

# Alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates


COPY --from=builder /go/src/github.com/seeleteam/monitor-api/build/monitor-api /monitor-api/

RUN cd  monitor-api

RUN chmod +x monitor-api

ENV MONITOR_CONFIG_FILE=/monitor-api/config/monitor.json

EXPOSE 9997

ENTRYPOINT [ "/monitor-api/monitor-api", "start", "-c=/monitor-api/config/app.conf" ]

# start monitor-api with your 'app.conf' file, this file must be external from a volume
# For example:
#   docker run -v <your app.conf path>:/monitor-api/config/app.conf:ro -it monitor-api