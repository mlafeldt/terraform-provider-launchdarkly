TEST ?= ./...

all: lint test testacc install

lint:
	go vet ./...
	golint ./...

test:
	go test $(TEST) -v $(TESTARGS)

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)

install:
	go install

install-deps:
	go get -d -t -v ./...
	go get github.com/golang/lint/golint
