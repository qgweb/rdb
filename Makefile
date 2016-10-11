GO := go

build:
	$(GO) build -o rdb server.go
clear:
	rm -f rdb