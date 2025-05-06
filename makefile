IMAGE_NAME ?= converter-routine-go-refactoring

.PHONY:
build-docker: build-go
	@echo "Building Docker image: $(IMAGE_NAME)"
	docker build -t $(IMAGE_NAME) .

.PHONY:
setup:
	@echo "Setting up dependencies..."
	# npm install --no-save swagger2openapi
	go install github.com/zxmfke/swagger2openapi3/cmd/swag2op@v1.0.1
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

.PHONY:
gen-docs: setup swag-format
	@echo "Generating Swagger 2.0 documentation..."
	swag2op init --openo ./docs/

.PHONY:
swag-format:
	@echo "Formatting Swag comments..."
	swag fmt

# not needed anymore, necessary when using the npm package for converting
# .PHONY:
# convert-swagger: setup
# 	@echo "Converting Swagger 2.0 to OpenAPI 3.0..."
# 	node_modules/.bin/swagger2openapi ./docs/swagger.json -o ./server/openapi.json; \

.PHONY:
build-go: gen-docs clean
	@echo "Building Go binary with embedded OpenAPI spec..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . 

.PHONY:
clean:
	@echo "Cleaning generated documentation..."
	rm -rf ./docs

.PHONY:
clean-all: clean
	rm -f ./openapi.json
	@echo "Removing npm modules..."
	rm -rf node_modules
