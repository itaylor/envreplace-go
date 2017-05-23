GOOS=linux

build:
	go build -o envreplace envreplace.go

test: build
	go test

install: test
	cp envreplace /usr/local/bin/envreplace
	chmod +x envreplace
