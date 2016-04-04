.PHONY: deps glide

deps:
	go get github.com/Masterminds/glide
	glide install
