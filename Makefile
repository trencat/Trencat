# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# DOCS parameters
DOCSDIR=$(CURDIR)/doc
DOCSFORMAT=html

.PHONY: all
all: test build

.PHONY: build build-train
build: build-train
build-train:
	$(GOBUILD) -v github.com/trencat/Trencat/train/

.PHONY: test test-update
test:
	$(GOTEST) -v github.com/trencat/Trencat/...
test-update:
	$(GOTEST) -v github.com/trencat/Trencat/... --update

.PHONY: run-logserver
run-logserver:
	docker-compose up -d logserver

.PHONY: shutdown
shutdown:
	docker-compose down

.PHONY: docs
docs:
	docker image build -t trencat_doc:poc $(DOCSDIR) && docker container run --rm -v $(DOCSDIR):/trencat_doc trencat_doc:poc sphinx-build -b $(DOCSFORMAT) /trencat_doc/source /trencat_doc/build/$(DOCSFORMAT)






#clean:
#	$(GOCLEAN) -cache
#	rm -f $(BINARY_NAME)
#	rm -f $(BINARY_UNIX)
#run:
#	$(GOBUILD) -o $(BINARY_NAME) -v ./...
#	./$(BINARY_NAME)
#deps:
#	$(GOGET) github.com/markbates/goth
#	$(GOGET) github.com/markbates/pop
