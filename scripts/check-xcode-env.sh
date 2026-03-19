#!/usr/bin/env bash
# Ensures the active Swift toolchain is full Xcode (not Command Line Tools only).
# Usage: ./scripts/check-xcode-env.sh
set -euo pipefail

dev_dir="${DEVELOPER_DIR:-$(xcode-select -p)}"

if [[ "$dev_dir" == *CommandLineTools* ]]; then
  echo "Prism must be built with Xcode.app, not Command Line Tools alone." >&2
  echo "" >&2
  echo "Fix:" >&2
  echo "  1. Install Xcode from the App Store (or Apple developer downloads)." >&2
  echo "  2. sudo xcode-select -s /Applications/Xcode.app/Contents/Developer" >&2
  echo "  3. sudo xcodebuild -license accept   # first launch only" >&2
  echo "  4. xcodebuild -runFirstLaunch        # optional, completes bundled components" >&2
  exit 1
fi

if [[ ! -x "$dev_dir/usr/bin/xcodebuild" ]]; then
  echo "xcodebuild not found under: $dev_dir" >&2
  exit 1
fi

echo "Xcode environment OK ($dev_dir)"
xcodebuild -version
