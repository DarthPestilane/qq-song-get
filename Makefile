.PHONY: all test clean

args = -ldflags='-s -w' -o qq-song-get

build:
	CGO_ENABLED=0 go build ${args}

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${args}

release:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=1.0.1 -X main.BuildTime=`date -u +\"%Y-%m-%dT%H:%M:%SZ\"`"

lint: # 代码风格检查
	golangci-lint run --concurrency=2

test:
	go test ./... -count=1

spec: lint test

tidy:
	go mod tidy
