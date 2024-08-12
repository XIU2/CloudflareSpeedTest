go env -w GOOS=linux
go build -ldflags "-w -s" -o dist/CloudflareSpeedTest .
go env -w GOOS=windows