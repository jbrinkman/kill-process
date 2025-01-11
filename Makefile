# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
BINARY_NAME=kp

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Clean the project
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Install dependencies
deps:
	$(GOGET) -u github.com/spf13/cobra/cobra

# Format the code
format:
	$(GOFMT) -w .

.PHONY: build clean test deps format