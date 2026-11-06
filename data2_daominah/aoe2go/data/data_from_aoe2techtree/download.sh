#!/bin/bash

# Downloads the latest data.json and strings.json from the aoe2techtree repo
# (its default branch is "master"), overwriting the files next to this script.
# Run "git diff --stat" afterwards to see what changed upstream.

set -euo pipefail

cd "$(dirname "$0")"

base_url="https://raw.githubusercontent.com/SiegeEngineers/aoe2techtree/master/data"

# download to a temp file first,
# so a failed download does not leave a truncated JSON behind
download() {
    local url="$1"
    local target="$2"
    curl --fail --silent --show-error --location "$url" -o "$target.tmp"
    mv -f "$target.tmp" "$target"
    echo "$(date --iso=s) downloaded $target ($(wc -c < "$target") bytes)"
}

download "$base_url/data.json" "data.json"
download "$base_url/locales/en/strings.json" "strings.json"
