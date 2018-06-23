prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/thisisaaronland/go-artisanal-integers; then rm -rf src/github.com/thisisaaronland/go-artisanal-integers; fi
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/client
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/engine
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/server
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/service
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/util
	cp *.go src/github.com/thisisaaronland/go-artisanal-integers/
	cp -r client src/github.com/thisisaaronland/go-artisanal-integers/
	cp -r engine src/github.com/thisisaaronland/go-artisanal-integers/
	cp -r server src/github.com/thisisaaronland/go-artisanal-integers/
	cp -r service src/github.com/thisisaaronland/go-artisanal-integers/
	cp -r util src/github.com/thisisaaronland/go-artisanal-integers/
	# if test -d vendor; then cp -r vendor/* src/; fi

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	# @GOPATH=$(shell pwd) go get "github.com/facebookgo/grace/gracehttp"

vendor-deps: rmdeps deps
	if test ! -d src; then mkdir src; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src
fmt:
	go fmt *.go
	go fmt client/*.go
	go fmt cmd/*.go
	go fmt engine/*.go
	go fmt server/*.go
	go fmt service/*.go
	go fmt util/*.go

bin:    self
	if test ! -d bin; then mkdir bin; fi
	@GOPATH=$(shell pwd) go build -o bin/int cmd/int.go
	@GOPATH=$(shell pwd) go build -o bin/intd-client cmd/intd-client.go
	@GOPATH=$(shell pwd) go build -o bin/intd-server cmd/intd-server.go

docker-build:
	docker build -t intd-server .