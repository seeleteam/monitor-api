@echo off 
goto comment
    Build the command lines and tests in Windows.
    Must install gcc tool before building.
:comment

echo on

go build -ldflags "-s -w" -o ./build/monitor-api.exe ./cmd/api
@echo "Done monitor-api building release"

pause
