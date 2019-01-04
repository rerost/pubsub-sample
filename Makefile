.PHONY: build
build:
	@go build -o bin/publisher publisher/main.go
	@go build -o bin/subscriber subscriber/main.go

.PHONY: clear
clear:
	rm -rf bin
