#!/bin/bash

# --------------------------------------------------
current_path=`pwd`
#compile and output files to
#target="./output/bin/monitor-api"
target="./bin/monitor-api"

# --------------------------------------------------
#get govendor
echo "get govendor"
cd $current_path
go get github.com/kardianos/govendor

# --------------------------------------------------
#build
echo "build"
if [ $@ = "debug" ]
then
    go build -o $target && echo "ok"
else
    go build -ldflags "-s -w" -o $target && echo "ok"
fi
mkdir -p ./output/bin
cp -rf $target ./output/bin/monitor-api 2>/dev/null
# --------------------------------------------------
