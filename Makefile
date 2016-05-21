build:
	go build -ldflags="-w" -o bin/ec2c

deps:
	go get github.com/Masterminds/glide
	glide install

.PHONY: build deps
