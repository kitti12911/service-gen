# syntax=docker/dockerfile:1.7
# check=skip=InvalidDefaultArgInFrom
ARG TOOLCHAIN_IMAGE
FROM ${TOOLCHAIN_IMAGE} AS builder

WORKDIR /src

# GOPRIVATE is baked in at build time so `go mod download` skips the public
# proxy/sumdb for modules hosted at the lib-path host. Empty when public.
ARG GOPRIVATE=___GOPRIVATE___
ENV GOPRIVATE=${GOPRIVATE}

COPY go.mod go.sum ./
# Optional buildx secret: pass `--secret id=netrc,src=$HOME/.netrc` (locally)
# or use the GitLab CI `.docker_job` template (mounts CI_JOB_TOKEN-based
# netrc). The mount is required=false so public-module builds keep working.
RUN --mount=type=secret,id=netrc,target=/root/.netrc,required=false \
	go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/___NAME___ ./cmd/server

FROM alpine:3.22@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/___NAME___ /app/___NAME___
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE 8080

ENTRYPOINT ["/app/___NAME___"]
