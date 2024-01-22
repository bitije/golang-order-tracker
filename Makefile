BINARY_NAME=ordertracker
.DEFAULT_GOAL := RUN

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux ./cmd/app/main.go
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows ./cmd/app/main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin ./cmd/app/main.go

config:
	docker compose up -d

run: 
	./$(BINARY_NAME)-linux

