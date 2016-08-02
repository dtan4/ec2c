FROM golang:1.6-onbuild

RUN make build
ENTRYPOINT ["/go/src/app/bin/ec2c"]
