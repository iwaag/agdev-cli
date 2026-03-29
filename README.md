# agdev

`agdev` is a CLI for agent-oriented AGDEV operations.

## What It Can Do

- print CLI version information
- read embedded instruction documents
- call authenticated backend operations such as `code mission`

## Install

Install the latest Linux release:

```bash
curl -fsSL https://raw.githubusercontent.com/iwaag/agdev-cli/main/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/iwaag/agdev-cli/main/install.sh | sh -s -- --version v0.1.0
```

## Basic Usage

Check the installed version:

```bash
agdev version
```

Read the common instruction text:

```bash
agdev code instruction common
```

Read a specific instruction version:

```bash
agdev code instruction common --version v1
```

Fetch an OpenAPI document and save it locally:

```bash
agdev util openapi http://127.0.0.1:8080
```

Keep only specific tagged endpoints:

```bash
agdev util openapi http://127.0.0.1:8080 --tags mission,auth -o .hints/backend.json
```

## Authentication

Log in once and save a local session:

```bash
agdev login --user <user_name>
```

If `--user` is omitted, `agdev` uses `KEYCLOAK_USER_NAME` when available. If that is also missing, it prompts for the username. Password entry is interactive and hidden.

`login` requires these environment variables:

- `KEYCLOAK_URL`
- `KEYCLOAK_REALM`
- `KEYCLOAK_CLIENT_ID`

Optional:

- `KEYCLOAK_USER_NAME`

After a successful login, `agdev` stores a local session and uses it automatically for authenticated commands.

Example:

```bash
export AGCODE_API_URL=http://127.0.0.1:8080
agdev login --user alice
agdev code mission mission-123
```

## Notes

- `code mission` requires `AGCODE_API_URL`.
- `code instruction common` reads static instruction text embedded in the binary and does not require login.
- Development, build, Docker, and release details are documented in `README_DEV.md`.
