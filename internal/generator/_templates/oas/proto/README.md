# proto

This OpenAPI service does not consume protobuf by default. If you add upstream
gRPC clients, drop `.proto` files here and wire a `buf.gen.yaml` like the
`grpc` pattern does.
