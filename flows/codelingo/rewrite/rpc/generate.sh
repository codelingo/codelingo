#!/usr/bin/env sh

# Install proto3 from source
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#
# Update protoc Go bindings via
#  go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

# (waigani) NOTE: there is a bug with the generated API. If you get this error message:
# ```cannot use _Query_Query_Handler (type func(interface {}, context.Context, func(interface {}) error) (interface {}, error)) as type grpc.methodHandler in field value```
# you need to add `interceptor grpc.UnaryServerInterceptor` as the last param of _Query_Query_Handler func sig.

# go
protoc $GOPATH/src/github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc/rewrite.proto -I=$GOPATH/src --go_out=plugins=grpc:$GOPATH/src
