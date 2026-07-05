# syntax=docker/dockerfile:1.7
# check=skip=InvalidDefaultArgInFrom
ARG TOOLCHAIN_IMAGE
FROM ${TOOLCHAIN_IMAGE} AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY buf.gen.yaml ./
COPY proto ./proto
COPY cmd ./cmd
COPY internal ./internal

RUN rm -rf gen/grpc \
	&& if [ -n "$(find proto -name '*.proto' 2>/dev/null)" ] \
		|| grep -qE '^[[:space:]]*-[[:space:]]*git_repo:' buf.gen.yaml 2>/dev/null; then \
		buf generate; \
	fi

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/___NAME___ ./cmd/server

FROM alpine:3.23@sha256:fd791d74b68913cbb027c6546007b3f0d3bc45125f797758156952bc2d6daf40

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/___NAME___ /app/___NAME___
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE 50051

ENTRYPOINT ["/app/___NAME___"]
