FROM alpine
CMD ["/bin/sh"]
ENV GO111MODULE=on GOBIN=/usr/local/bin GOPATH=/go
RUN apk add --no-cache ca-certificates go openssl nodejs npm
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogo
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogofast
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
RUN go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
RUN go get -u github.com/bufbuild/buf github.com/bufbuild/buf/cmd/buf github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking github.com/bufbuild/buf/cmd/protoc-gen-buf-lint
RUN npm install -g swagger-combine
