FROM alpine:3.4
MAINTAINER Daisuke Fujita <dtanshi45@gmail.com> (@dtan4)

RUN apk add --no-cache --update ca-certificates

COPY bin/ec2c /ec2c

ENTRYPOINT ["/ec2c"]
