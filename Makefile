build: fmt
	CGO_ENABLED=0 GOOS=linux go build -v -o ghcs

fmt:
	go fmt ./...