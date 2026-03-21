# agdev

`agdev` is a Go CLI skeleton for agent-oriented tooling.

The primary distribution model is a normal CLI binary. Docker support exists as an optional packaging and runtime path.

## Current Scope

- Cobra-based command tree
- Environment-driven configuration
- Text or JSON output mode
- Installable local CLI binary
- Optional Docker-based execution path
- Versioned static instruction text for agents

## Local CLI Usage

Install the latest Linux release:

```bash
curl -fsSL https://raw.githubusercontent.com/iwaag/agdev-cli/main/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/iwaag/agdev-cli/main/install.sh | sh -s -- --version v0.1.0
```

Build a local binary:

```bash
make build
./bin/agdev version
./bin/agdev code instruction common
```

Install into your Go bin directory:

```bash
make install
agdev version
agdev code instruction common
```

Direct execution during development:

```bash
go run . version
go run . code instruction common
go run . code instruction common --version latest --json
```

## Environment Variables

- `AGDEV_OUTPUT_JSON`
- `AGDEV_LOG_LEVEL`

## Docker

Build a standalone CLI image:

```bash
docker build -t agdev .
```

Build `agdev` from this GitHub repository inside another Dockerfile:

```dockerfile
FROM golang:1.23-alpine AS agdev-builder

RUN apk add --no-cache git

WORKDIR /src
RUN git clone https://github.com/your-org/agdev-cli.git .
RUN go build -o /out/agdev .

FROM alpine:3.20
COPY --from=agdev-builder /out/agdev /usr/local/bin/agdev

ENTRYPOINT ["/usr/local/bin/agdev"]
```

Publish a GitHub release build:

```bash
./scripts/release.sh patch --push
```

The release script creates the next `v*` tag from the latest existing tag. Pushing a tag that matches `v*` triggers the GitHub Actions release workflow, which builds `agdev` for Linux amd64 and uploads the archive and checksum to GitHub Releases.

You can also open the `Release` workflow in GitHub Actions and run it manually with an existing tag such as `v0.1.0` to rebuild and replace the release assets for that tag.

You can also create a tag without pushing it yet:

```bash
./scripts/release.sh patch
git push origin <new-tag>
```

Run the CLI from that image:

```bash
docker run --rm agdev version
docker run --rm agdev code instruction common
```

Copy the binary into another service image:

```dockerfile
FROM agdev:latest AS agdev-cli

FROM your-service-image
COPY --from=agdev-cli /usr/local/bin/agdev /usr/local/bin/agdev
```

## Notes

The `code instruction common` command reads versioned static instruction text embedded in the CLI binary.
