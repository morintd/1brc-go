install:
	go get ./...
build:
	go build -o bin/application main.go
test:
	"$(CURDIR)/scripts/test.sh"
install_tools:
	"$(CURDIR)/scripts/install_tools.sh"

.NOTPARALLEL:

.PHONY: install build test install_tools 