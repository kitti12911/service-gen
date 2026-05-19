# **_NAME_**

Background worker service scaffolded by
[`service-gen`](https://github.com/kitti12911/service-gen).

## What's included

- NATS / JetStream consumer using `lib-async` with typed JSON job decoding.
- At-least-once delivery with ack / nack handling.
- Configuration via `config.yml` / environment variables (`lib-util/v3`).
- Structured logging, OpenTelemetry tracing, and Pyroscope profiling, with
  trace propagation from publisher messages to the worker handler.
- A `internal/worker` package with an example job payload and handler.
- GitHub Actions and/or GitLab CI, Renovate, CODEOWNERS, golangci-lint,
  markdownlint, Prettier, `.air.toml`, and a multi-stage `Dockerfile`.

## Getting started

```sh
cp config.example.yml config.yml   # then edit values, including the NATS topic
make run
```

A local NATS broker is required at runtime.

<!-- IF_GOPRIVATE -->
## Private Go modules

This project's libraries are hosted at `___LIB_PATH___`. `GOPRIVATE=___GOPRIVATE___`
is preset in the Dockerfile and GitLab CI variables so the Go toolchain skips
the public proxy/sumdb for that host.

### Cloning the lib repos (one-time local setup)

Pick one auth method.

**SSH (recommended for humans):**

```sh
# Add your SSH key to GitLab → User Settings → SSH Keys, then:
git clone git@___GOPRIVATE___:___LIB_NAMESPACE___/lib-util.git

# Tell git to rewrite https://___GOPRIVATE___/ to SSH so `go mod tidy` uses it.
git config --global url."git@___GOPRIVATE___:".insteadOf "https://___GOPRIVATE___/"
```

**Personal access token (works for headless / CI-like setups):**

```sh
# Create a token at GitLab → User Settings → Access Tokens with scope `read_api`
# and `read_repository`, then:
printf 'machine ___GOPRIVATE___ login <gitlab-username> password <token>\n' \
    >> ~/.netrc
chmod 600 ~/.netrc

git clone https://___GOPRIVATE___/___LIB_NAMESPACE___/lib-util.git
```

Then `make tidy` / `make run` resolve `___LIB_PATH___/lib-*` modules normally.

### CI and Docker builds

- **GitLab CI**: `CI_JOB_TOKEN` is used automatically — the top-level
  `default.before_script` writes `~/.netrc` before any Go command runs.
- **Docker build**: pass `--secret id=netrc,src=$HOME/.netrc` to `docker build`
  locally. CI exports `BUILD_SECRETS_NETRC` from `.docker_job` and
  `scripts/ci/build-image.sh` forwards it as a buildx secret. Example:

  ```sh
  docker buildx build \
      --secret id=netrc,src=$HOME/.netrc \
      --build-arg TOOLCHAIN_IMAGE=<your-toolchain-image> \
      -t ___NAME___ .
  ```
<!-- END_GOPRIVATE -->

## Message shape

```json
{
    "id": "job-1",
    "type": "debug.print",
    "payload": { "message": "hello" }
}
```

## Delivery behavior

With JetStream enabled the worker uses at-least-once delivery: a handler error
nacks the message and the broker can redeliver it, so handlers should be
idempotent.

## Common commands

| Command       | Description                               |
| ------------- | ----------------------------------------- |
| `make run`    | Start the worker locally                  |
| `make air`    | Run with live reload                      |
| `make test`   | Run tests with the race detector          |
| `make lint`   | Run Go and Markdown linting               |
| `make format` | Format Go, Markdown, YAML, JSON           |
| `make cov`    | Generate and open an HTML coverage report |

Implement your job logic in `internal/worker/handler.go`.

## Deployment

A starter Helm chart is included in [`helm/`](helm/) so the worker is
deployable out of the box. Deployment state should not live in the service
repository long-term — move the chart to your `helm-sandbox` repository (or
your preferred chart repository), then delete `helm/` from here:

```sh
mv helm /path/to/helm-sandbox/charts/___NAME___
```

Keep app-specific values (image, env, scaling, config) in the chart
repository, not in the service repo.
