.PHONY: deps

deps:
	go get github.com/Masterminds/glide
	glide install
