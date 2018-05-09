# Go builder container
FROM golang:alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go/src/github.com/seeleteam/monitor-api

WORKDIR /go/src/github.com/seeleteam/monitor-api
RUN cd /go/src/github.com/seeleteam/monitor-api

RUN go build -ldflags "-s -w" -o ./run-monitor-api .

# Alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates


COPY --from=builder /go/src/github.com/seeleteam/monitor-api/run-monitor-api /monitor-api/

RUN chmod +x /monitor-api/run-monitor-api

EXPOSE 9997

ENTRYPOINT [ "/monitor-api/run-monitor-api" ]

# start monitor-api with your 'app.conf' file, this file must be external from a volume
# For example:
#   docker run -v <your app.conf path>:/monitor-api/conf/app.conf:ro -it monitor-api