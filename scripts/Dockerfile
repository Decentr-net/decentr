FROM golang:alpine AS build-env

WORKDIR /go/scr/github.com/decentr/decentr/

RUN apk add --update make git

COPY . .

RUN make install

FROM alpine:edge

RUN apk add --update ca-certificates

COPY --from=build-env /go/bin/decentrd /usr/bin/decentrd
COPY --from=build-env /go/bin/decentrcli /usr/bin/decentrcli

EXPOSE 26657
EXPOSE 26656
EXPOSE 1317

CMD ["decentrd", "start"]