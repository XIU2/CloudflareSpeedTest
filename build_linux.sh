#!/bin/bash
go env -w GOOS=linux
go build -ldflags "-w -s" -o dist/CloudflareSpeedTest .