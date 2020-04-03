build: fmt
	CGO_ENABLED=0 GOOS=linux go build -v -o ghcs

fmt:
	go fmt ./...

dbuild: fmt
	docker build --tag caladreas/ghcs:latest .

dpush: dbuild
	docker push caladreas/ghcs:latest