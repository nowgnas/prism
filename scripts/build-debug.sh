#!/usr/bin/env bash
# Generate project (xcodegen) and Debug build. Run from repo root.
set -euo pipefail
cd "$(dirname "$0")/.."

./scripts/check-xcode-env.sh
command -v xcodegen >/dev/null || {
  echo "Install xcodegen: brew install xcodegen" >&2
  exit 1
}

xcodegen generate
xcodebuild \
  -project Prism.xcodeproj \
  -scheme Prism \
  -configuration Debug \
  build \
  ARCHS="x86_64 arm64" \
  CODE_SIGNING_REQUIRED=NO \
  CODE_SIGNING_ALLOWED=NO
