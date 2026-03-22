# agdev Developer Guide

This document is for development, packaging, and operational usage.

## Local Development

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
```

## Auth Model

There are two intended auth patterns.

### Interactive User Flow

Use `agdev login` to create a saved local session.

Required environment variables:

- `KEYCLOAK_URL`
- `KEYCLOAK_REALM`
- `KEYCLOAK_CLIENT_ID`

Optional:

- `KEYCLOAK_USER_NAME`

The saved session is used automatically by authenticated commands unless `--token` is supplied.

### Backend Service-To-Service Flow

For backend integration or automated service-to-service calls, specify `--token` on each command. This is a system-side usage pattern, not the normal end-user flow.

Example:

```bash
AGCODE_API_URL=http://127.0.0.1:8080 \
agdev code mission mission-123 --token "$ACCESS_TOKEN"
```

This is the expected pattern for non-interactive callers. Do not rely on a local login session in backend service execution paths.

### Resolution Order

Authenticated command resolution is:

1. `--token`
2. saved local session
3. auth failure

## Docker

Build a standalone CLI image:

```bash
docker build -t agdev .
```

Run the CLI from that image:

```bash
docker run --rm agdev version
docker run --rm agdev code instruction common
```

Run an authenticated backend call with an explicit token:

```bash
docker run --rm \
  -e AGCODE_API_URL=http://host.docker.internal:8080 \
  agdev code mission mission-123 --token "$ACCESS_TOKEN"
```

Build `agdev` from this repository inside another Dockerfile:

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

Copy the binary into another service image:

```dockerfile
FROM agdev:latest AS agdev-cli

FROM your-service-image
COPY --from=agdev-cli /usr/local/bin/agdev /usr/local/bin/agdev
```

## Release

Publish a GitHub release build:

```bash
./scripts/release.sh patch --push
```

The release script creates the next `v*` tag from the latest existing tag. Pushing a matching tag triggers the GitHub Actions release workflow, which builds `agdev` for Linux amd64 and uploads the archive and checksum to GitHub Releases.

Create a tag without pushing it:

```bash
./scripts/release.sh patch
git push origin <new-tag>
```

You can also open the `Release` workflow in GitHub Actions and run it manually with an existing tag such as `v0.1.0` to rebuild and replace the release assets for that tag.
