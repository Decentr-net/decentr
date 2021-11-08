ARG ALPINE_VERSION=3.10
ARG GOLANG_VERSION=1.16.9

FROM golang:${GOLANG_VERSION} AS build-env
WORKDIR /go/src/github.com/Decentr-net/decentr/
COPY . .
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --update ca-certificates libc6-compat
COPY --from=build-env /go/src/github.com/Decentr-net/decentr/build/decentr-linux-amd64 /usr/bin/decentr

EXPOSE 26657
EXPOSE 26656
EXPOSE 9090

CMD ["decentr", "start"]