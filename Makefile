TEST ?= ./...

all: lint test testacc build

install-deps:
	go get -d -t -v ./...
	go get github.com/golang/lint/golint

lint:
	go vet ./...
	golint ./...

test:
	go test $(TEST) -v $(TESTARGS)

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)

build:
	go build -o terraform-provider-launchdarkly

example: build
	terraform init
	terraform apply
