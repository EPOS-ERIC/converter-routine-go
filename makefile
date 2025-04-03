.PHONY: gen-docs convert-swagger clean setup build-go build-docker clean-all

# Build Docker image 
build-docker: build-go
	@echo "Building Docker image: $(IMAGE_NAME)"
	$(eval IMAGE_NAME := converter-routine-go-refactoring)
	docker build -t $(IMAGE_NAME) .

# Install necessary npm and go packages
setup:
	@echo "Setting up dependencies..."
	npm install --no-save swagger2openapi
	go get -d -v ./...
	go install -v ./...
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger 2.0 docs using swag
gen-docs:
	@echo "Generating Swagger 2.0 documentation..."
	swag init -g main.go

# Convert Swagger 2.0 to OpenAPI 3.0 using npm package
convert-swagger:
	@echo "Converting Swagger 2.0 to OpenAPI 3.0..."
	@if [ -f "./docs/swagger.json" ]; then \
		node_modules/.bin/swagger2openapi ./docs/swagger.json -o ./openapi.json; \
		echo "Conversion complete. OpenAPI 3.0 specs saved to ./openapi.json"; \
	else \
		echo "Error: Swagger 2.0 specs not found at ./docs/swagger.json"; \
		exit 1; \
	fi

# Build the Go application (with embedded OpenAPI file)
build-go: gen-docs convert-swagger
	@echo "Building Go binary with embedded OpenAPI spec..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . 

# Clean generated documentation
clean:
	@echo "Cleaning generated documentation..."
	rm -f ./docs/swagger.json ./docs/swagger.yaml ./openapi.json
	rm -f ./docs/docs.go

clean-all: clean
	@echo "Removing npm modules..."
	rm -rf node_modules
