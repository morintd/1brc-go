install:
	go get ./...
build:
	go build -o bin/application main.go
build_release:
	CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o bin/application main.go
test:
	"$(CURDIR)/scripts/test.sh"
install_tools:
	"$(CURDIR)/scripts/install_tools.sh"

.NOTPARALLEL:

.PHONY: install build build_release test install_tools 