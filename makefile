all:
	CGO_ENABLED=0 go build -v
start:
	./whiskercat
clean:
	go fmt ./...