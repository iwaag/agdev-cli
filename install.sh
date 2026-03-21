#!/usr/bin/env bash

set -eu

repo="iwaag/agdev-cli"
asset_name="agdev_linux_amd64.tar.gz"
version="latest"
bin_dir=""

usage() {
  cat <<'EOF'
Usage:
  install.sh [--version <version>] [--bin-dir <dir>]

Options:
  --version <version>  Install a specific release tag such as v0.1.0.
                       Defaults to latest.
  --bin-dir <dir>      Install directory for the agdev binary.
                       Defaults to /usr/local/bin when writable,
                       otherwise ~/.local/bin.
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --version)
      if [ $# -lt 2 ]; then
        usage
        exit 1
      fi
      version="$2"
      shift 2
      ;;
    --bin-dir)
      if [ $# -lt 2 ]; then
        usage
        exit 1
      fi
      bin_dir="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [ "$(uname -s)" != "Linux" ]; then
  echo "install.sh currently supports Linux only" >&2
  exit 1
fi

arch="$(uname -m)"
if [ "$arch" != "x86_64" ] && [ "$arch" != "amd64" ]; then
  echo "install.sh currently supports linux amd64 only" >&2
  exit 1
fi

if [ -z "$bin_dir" ]; then
  if [ -w /usr/local/bin ]; then
    bin_dir="/usr/local/bin"
  else
    bin_dir="${HOME}/.local/bin"
  fi
fi

mkdir -p "$bin_dir"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT INT TERM

archive_path="${tmpdir}/${asset_name}"
checksum_path="${archive_path}.sha256"

if command -v curl >/dev/null 2>&1; then
  download() {
    curl -fsSL "$1" -o "$2"
  }
elif command -v wget >/dev/null 2>&1; then
  download() {
    wget -qO "$2" "$1"
  }
else
  echo "curl or wget is required" >&2
  exit 1
fi

if [ "$version" = "latest" ]; then
  base_url="https://github.com/${repo}/releases/latest/download"
else
  base_url="https://github.com/${repo}/releases/download/${version}"
fi

download "${base_url}/${asset_name}" "$archive_path"
download "${base_url}/${asset_name}.sha256" "$checksum_path"

(
  cd "$tmpdir"
  sha256sum -c "$(basename "$checksum_path")"
)

tar -xzf "$archive_path" -C "$tmpdir"
install -m 0755 "${tmpdir}/agdev" "${bin_dir}/agdev"

echo "installed agdev to ${bin_dir}/agdev"
