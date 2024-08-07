run:
	go run main.go

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o bin/dbdl-amd64-linux

deploy:
	./deploy.sh
