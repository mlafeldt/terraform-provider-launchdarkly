TEST ?= ./...

all: lint test testacc install

lint:
	go vet $(TEST)
	golint $(TEST)

test:
	go test -i -v $(TEST)
	go test $(TEST) -v $(TESTARGS)

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)

install:
	go install -v
