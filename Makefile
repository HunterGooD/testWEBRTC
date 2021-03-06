APP?=testwebrtc
GOOS?=linux
GOARCH?=amd64

.PHONY: clean, build

build: clean
	cd cmd/${APP} && \
	go build -o ../../dist/${APP}

clean:
	@rm -rf ./dist

dev:
	PORT=8080 go run cmd/testwebrtc/main.go
