# ==============================================================================
# Main commands

run:
	@go run ./cmd/email_service/main.go

build: 
	@go build ./cmd/email_service/main.go

tests: 
	go test -v ./test/...

clean:
	go mod tidy && go fmt ./...

lint:
	golangci-lint run ./...

# ==============================================================================
# Docker
compose_run:
	@docker-compose up --force-recreate -f ./deployments/docker-compose.yml 
