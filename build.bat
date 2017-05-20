@echo off
echo Building
go get
go build -o codegen.exe -ldflags "-s -w"