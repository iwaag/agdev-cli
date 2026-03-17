# agdev

`agdev` is a Go CLI skeleton for agent-driven execution against AGDEV backends.

## Current Scope

- Cobra-based command tree
- Environment-driven configuration
- Text or JSON output mode
- Docker-first build and execution path
- Stub commands for image and video generation

## Commands

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

Build:

```bash
docker build -t agdev .
```

Run:

```bash
docker run --rm agdev version
docker run --rm agdev image generate input.png "describe the edit"
docker run --rm agdev video generate first.png last.png --json
```

## Notes

The current `generate` commands return stub responses only. The backend REST and socket integration points are prepared but not implemented yet.
