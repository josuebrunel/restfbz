.PHONY: all test coverage

clean:
	go clean -i -x

get:
	go get -u ./...

build:
	go build ./...

install:
	go install ./...

lint:
	golint ./...

test:
	go test ./... -v -p=1 -count=1 -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt

all: get build install
