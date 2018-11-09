TEST ?= ./...

all: lint test testacc build

install-deps:
	go get -d -t -v ./...
	go get github.com/golang/lint/golint

lint:
	go vet ./...
	golint -set_exit_status ./...

test:
	go test $(TEST) -v $(TESTARGS)

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS)

build: install-deps
	go build -o terraform-provider-launchdarkly

install: build
	install -d -m 755 ~/.terraform.d/plugins
	install terraform-provider-launchdarkly ~/.terraform.d/plugins

apply: build
	terraform init
	terraform $@

destroy: build
	terraform destroy
	terraform $@
