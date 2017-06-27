HANDLER ?= golambda
PACKAGE ?= $(HANDLER)
GOPATH  ?= $(HOME)/go

WORKDIR = $(CURDIR:$(GOPATH)%=/go%)
ifeq ($(WORKDIR),$(CURDIR))
	WORKDIR = /build
endif

all: test build pack perm

.PHONY: all

docker:
	@docker run --rm                                                             \
	  -e HANDLER=$(HANDLER)                                                      \
	  -e PACKAGE=$(PACKAGE)                                                      \
	  -v $(GOPATH):/go                                                           \
	  -v $(CURDIR):/build                                                        \
	  -w $(WORKDIR)                                                              \
	  pcrofts/bg-lambda-deploy:latest make all

.PHONY: docker

test:
	@go test `glide novendor` -cover

.PHONY: test

build:
	@go build -buildmode=plugin -ldflags='-w -s' -o $(HANDLER).so

.PHONY: build

pack:
	@pack $(HANDLER) $(HANDLER).so $(PACKAGE).zip

.PHONY: pack

perm:
	@chown $(shell stat -c '%u:%g' .) $(HANDLER).so $(PACKAGE).zip

.PHONY: perm

deploy:
	@aws lambda update-function-code --function-name golambda --zip-file file://$(PACKAGE).zip

.PHONY: deploy

clean:
	@rm -rf $(HANDLER).so $(PACKAGE).zip

.PHONY: clean