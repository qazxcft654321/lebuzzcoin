PROJECT_NAME=lebuzzcoin
GOOS=linux
BINARY_NAME=app
PKG_LIST := $(shell go list $(PROJECT_NAME)/...)

build: ## Build the binary file
	CGO_ENABLED=0 GOOS=$(GOOS) go build -a -installsuffix cgo -o $(BINARY_NAME) .

run: ## Run entire project
	@go run .

format: ## Format the files
	@go fmt ./...
	@go vet ./...

test: ## Run tests
	@go test -v $(PKG_LIST)

cover: ## Run tests coverage
	@go test -cover $(PKG_LIST)

loc: ## Print LOC of all .go files
	@find . -name '*.go' | xargs wc -l

help: ## Display the help screen
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
