#!/usr/bin/env bash

set -eu

usage() {
  cat <<'EOF'
Usage:
  scripts/release.sh patch [--push]
  scripts/release.sh minor [--push]
  scripts/release.sh major [--push]

Creates the next v* git tag from the latest existing tag.
If no matching tag exists yet, the first tag is v0.1.0.
Use --push to push the new tag to origin after creating it.
EOF
}

if [ $# -lt 1 ] || [ $# -gt 2 ]; then
  usage
  exit 1
fi

bump="$1"
push_tag="false"

if [ $# -eq 2 ]; then
  if [ "$2" != "--push" ]; then
    usage
    exit 1
  fi
  push_tag="true"
fi

case "$bump" in
  patch|minor|major)
    ;;
  *)
    usage
    exit 1
    ;;
esac

latest_tag="$(git tag --list 'v*' --sort=-version:refname | head -n1)"

if [ -z "$latest_tag" ]; then
  next_tag="v0.1.0"
else
  version="${latest_tag#v}"
  old_ifs="${IFS}"
  IFS=.
  set -- $version
  IFS="${old_ifs}"

  major="${1:-0}"
  minor="${2:-0}"
  patch="${3:-0}"

  case "$bump" in
    patch)
      patch=$((patch + 1))
      ;;
    minor)
      minor=$((minor + 1))
      patch=0
      ;;
    major)
      major=$((major + 1))
      minor=0
      patch=0
      ;;
  esac

  next_tag="v${major}.${minor}.${patch}"
fi

if git rev-parse --verify --quiet "$next_tag" >/dev/null; then
  echo "tag already exists: $next_tag" >&2
  exit 1
fi

git tag "$next_tag"
echo "created tag: $next_tag"

if [ "$push_tag" = "true" ]; then
  git push origin "$next_tag"
  echo "pushed tag: $next_tag"
fi
