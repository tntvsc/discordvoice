build:
	go build -o bin/main main.go

run:
	go run main.go

nodemon:
	nodemon --exec go run main.go --signal SIGTERM