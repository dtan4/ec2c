# ec2c

[![Build Status](https://travis-ci.org/dtan4/ec2c.svg?branch=master)](https://travis-ci.org/dtan4/ec2c)

Simple AWS EC2 CLI

## Usage

### `ec2c cancel`

Cancel the specified EC2 Spot Instance requests

### `ec2c launch`

Launch new EC2 instance

### `ec2c list`

List EC2 instances

### `ec2c list_requests`

List Spot Instance requests

### `ec2c request`

Request new Spot Instances

### `ec2c tag`

Tagging to EC2 instance

### `ec2c untag`

Delete tag from the specified EC2 instance

### `ec2c terminate`

Terminate the specified EC2 instance

## Install

### Precompiled binary

Precompiled binaries are available at [GitHub Releases](https://github.com/dtan4/ec2c/releases).

### From source

To install, use `go get`:

```bash
$ go get -d github.com/dtan4/ec2c
```

## Contribution

1. Fork ([https://github.com/dtan4/ec2c/fork](https://github.com/dtan4/ec2c/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[dtan4](https://github.com/dtan4)

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
