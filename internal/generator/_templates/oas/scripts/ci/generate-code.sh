#!/usr/bin/env sh
set -eu

# The generated OpenAPI service has no protobuf or patch codegen by default.
# This is a structural no-op so the CI generate stage stays uniform with the
# other service patterns. Add real generation here if you introduce upstream
# gRPC clients or PATCH mappers.

mkdir -p gen
# Non-hidden marker so upload-artifact (which excludes hidden files by
# default) still ships a non-empty "generated-code" artifact.
: >gen/keep

echo "no code generation required for the oas pattern"
