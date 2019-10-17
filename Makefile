args = -ldflags='-s -w' -o qq-song-get

build:
	CGO_ENABLED=0 go build ${args}

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${args}
