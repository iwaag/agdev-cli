# agdev

`agdev` is a Go CLI skeleton for agent-driven execution against AGDEV backends.

The primary distribution model is a normal CLI binary. Docker support exists as an optional packaging and runtime path.

## Current Scope

- Cobra-based command tree
- Environment-driven configuration
- Text or JSON output mode
- Installable local CLI binary
- Optional Docker-based execution path
- Stub commands for image and video generation

## Local CLI Usage

Build a local binary:

```bash
make build
./bin/agdev version
```

Install into your Go bin directory:

```bash
make install
agdev version
```

Direct execution during development:

```bash
go run . version
go run . image generate input.png "describe the edit"
go run . video generate first.png last.png --json
```

## Environment Variables

- `AGDEV_API_BASE_URL`
- `AGDEV_AUTH_TOKEN`
- `AGDEV_REQUEST_TIMEOUT`
- `AGDEV_OUTPUT_JSON`
- `AGDEV_LOG_LEVEL`

## Docker

Build a standalone CLI image:

```bash
docker build -t agdev .
```

Run the CLI from that image:

```bash
docker run --rm agdev version
docker run --rm agdev image generate input.png "describe the edit"
docker run --rm agdev video generate first.png last.png --json
```

Copy the binary into another service image:

```dockerfile
FROM agdev:latest AS agdev-cli

FROM your-service-image
COPY --from=agdev-cli /usr/local/bin/agdev /usr/local/bin/agdev
```

## Notes

The current `generate` commands return stub responses only. The backend REST and socket integration points are prepared but not implemented yet.
