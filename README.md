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
