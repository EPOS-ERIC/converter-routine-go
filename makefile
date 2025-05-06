IMAGE_NAME ?= converter-routine-go-refactoring

.PHONY:
build-docker: build-go
	@echo "Building Docker image: $(IMAGE_NAME)"
	docker build -t $(IMAGE_NAME) .

.PHONY:
setup:
	@echo "Setting up dependencies..."
	go install github.com/zxmfke/swagger2openapi3/cmd/swag2op@v1.0.1
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

.PHONY:
gen-docs: setup swag-format
	@echo "Generating Swagger 2.0 documentation..."
	swag2op init --openo ./docs/
	cp ./docs/swagger.json ./openapi.json

.PHONY:
swag-format:
	@echo "Formatting Swag comments..."
	swag fmt

.PHONY:
build-go: gen-docs
	@echo "Building Go binary with embedded OpenAPI spec..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . 

.PHONY:
clean:
	@echo "Cleaning generated documentation..."
	rm -rf ./docs
	rm -f ./openapi.json
