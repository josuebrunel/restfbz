.PHONY: all test coverage

clean:
	go clean -i -x

get:
	go get -u ./...

build:
	go build -o bin/restfbz cmd/restfbz/main.go

install:
	go install ./...

lint:
	golint ./...

test:
	go test ./... -v -p=1 -count=1 -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt

docker-build:
	docker-compose up --build -d

docker-test: docker-build
	docker exec -it restfbz_web make test

docker-run: docker-build
	docker exec -it restfbz_web go run cmd/restfbz/main.go 8999

all: get build install
