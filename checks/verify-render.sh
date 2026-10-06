#!/usr/bin/env bash
set -euo pipefail
# Source checks never install tools or execute template-declared commands.
cd "$(dirname "$0")"
GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local go test -count=1 ./...
