run:
	go run main.go

test:
	go test -v ./...

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o bin/sshdbd-amd64-linux
