.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO = go
BINARY_DIR=bin
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PACKAGES_NOVENDOR = $(shell glide novendor)

BUILD_DEPS:= github.com/alecthomas/gometalinter

.PHONY: vendor test build

help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

vendor:
	$(GO) get -v github.com/Masterminds/glide
	cd $(GOPATH)/src/github.com/Masterminds/glide && git checkout tags/v0.12.3 && go install && cd -
	glide install

get-build-deps:
	$(GO) get $(BUILD_DEPS)
	gometalinter --install

test:
	./gotest.sh

checkstyle: get-build-deps
	gometalinter --vendor ./... --fast --disable=gas --disable=errcheck --disable=gotype --deadline 10m

fmt:
	gofmt -l -w -s ${GOFILES_NOVENDOR}

# Builds the project
build: checkstyle test
	$(GO) build $(PACKAGES_NOVENDOR)

rewrite-import-paths:
	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"github.com/reportportal/commons-go%\"gopkg.in/reportportal/commons-go.v1%g '{}' \;

restore-import-paths:
	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"github.com/reportportal/commons-go%\"gopkg.in/reportportal/commons-go.v1%g '{}' \;

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi

release:
	git checkout -b ${v}
	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"github.com/reportportal/commons-go%\"gopkg.in/reportportal/commons-go.v1%g '{}' \;
	git add .
	git commit -m "rewrite path"
	git push origin ${v}
	git tag -a ${v} -m "creating tag ${v}"
	git push origin ${v}
	git checkout master
	git branch -D ${v}
	git push origin --delete "origin/${v}"
