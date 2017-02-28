# ec2c

[![Build Status](https://travis-ci.org/dtan4/ec2c.svg?branch=master)](https://travis-ci.org/dtan4/ec2c)

Simple AWS EC2 CLI

## Installation

### Using Homebrew (OS X only)

Formula is available at [dtan4/homebrew-dtan4](https://github.com/dtan4/homebrew-dtan4).

```bash
$ brew tap dtan4/dtan4
$ brew install ec2c
```

### Precompiled binary

Precompiled binaries are available at [GitHub Releases](https://github.com/dtan4/ec2c/releases).

### From source

```bash
$ go get -d github.com/dtan4/ec2c
$ cd $GOPATH/src/github.com/dtan4/ec2c
$ make
$ make instatll
```

## Usage

### Configuring credentials

Before using ec2c, ensure that you have credentials in `~/.aws/credentials`, which might look like:

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

Alternatively, you can set the following environment variables:

```
export AWS_ACCESS_KEY_ID=AKID1234567890
export AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

### Configuring the region

In addition to credentials, you must specify the region by the environmental variable:

```
export AWS_REGION=ap-northeast-1
```

Alternatively, you can enable support for the shared config `~/.aws/config`.
See [the aws-sdk-go documentation](https://github.com/aws/aws-sdk-go#configuring-credentials) for further information.

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
