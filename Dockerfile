FROM golang:alpine AS builder
COPY . /build-src
ENV GOPROXY='https://goproxy.io'
ENV GOSUMDB='off'
RUN cd /build-src && CGO_ENABLED=0 go build -v -ldflags='-s -w' -o /qq-song-get

FROM alpine:latest
COPY --from=builder /qq-song-get /usr/local/bin/
ENTRYPOINT ["qq-song-get"]
