prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/thisisaaronland/go-artisanal-integers; then rm -rf src/github.com/thisisaaronland/go-artisanal-integers; fi
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/engine
	mkdir -p src/github.com/thisisaaronland/go-artisanal-integers/service
	cp *.go src/github.com/thisisaaronland/go-artisanal-integers/
	cp engine/*.go src/github.com/thisisaaronland/go-artisanal-integers/engine/
	cp service/*.go src/github.com/thisisaaronland/go-artisanal-integers/service/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(shell pwd) go get "github.com/go-sql-driver/mysql"
	@GOPATH=$(shell pwd) go get "github.com/garyburd/redigo/redis"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go
	go fmt engine/*.go
	go fmt service/*.go

bin:    self
	@GOPATH=$(shell pwd) go build -o bin/int cmd/int.go
