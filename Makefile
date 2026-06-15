.PHONY: build install test tidy clean

build:
	go build -o bin/hr .

install:
	go install .

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin
