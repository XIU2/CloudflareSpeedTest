#!/bin/bash
go env -w GOOS=windows
go build -ldflags "-w -s" -o dist/CloudflareSpeedTest.exe .
go env -w GOOS=linux