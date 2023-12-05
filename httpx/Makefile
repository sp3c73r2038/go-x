default: test

test:
	go test -p=1 -v .

build:
	go build -o bin/http http.go
	go build -o bin/debug-server debug-server.go

run: build
	bin/http
