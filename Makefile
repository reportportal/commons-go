.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO = go
BINARY_DIR=bin
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: test build

help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

get-build-deps:
	# installs golangci-lint
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test:
	./gotest.sh

checkstyle:
	golangci-lint run --deadline=60m ./...


fmt:
	gofmt -l -w -s ${GOFILES_NOVENDOR}

# Builds the project
build: test
	$(GO) build ./...

rewrite-import-paths:
	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"github.com/reportportal/commons-go%\"gopkg.in/reportportal/commons-go.v5%g '{}' \;

restore-import-paths:
	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"gopkg.in/reportportal/commons-go.v5%\"github.com/reportportal/commons-go%g '{}' \;

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi

release:
#	git checkout -b temp-${v}
#	find . -not -path "./vendor/*" -name '*.go' -type f -execdir sed -i '' s%\"github.com/reportportal/commons-go%\"gopkg.in/reportportal/commons-go.v5%g '{}' \;
#	git add .
#	git status
#   git commit -m "rewrite import paths"
#	git push --set-upstream origin temp-${v}
	git tag -a ${v} -m "creating tag ${v}"
	git push origin "refs/tags/${v}"
	git checkout master
# 	git branch -D temp-${v}
# 	git push origin --delete temp-${v}
