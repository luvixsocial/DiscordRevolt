all:
	CGO_ENABLED=0 go build -v
start:
	./libdozina
clean:
	go fmt ./...