#!/usr/bin/env bash
set -euo pipefail

# -------------------------
# Ikaitla CLI Installer
# -------------------------

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Error: '$1' is required but not found in PATH." >&2
    exit 1
  }
}

prompt_default() {
  local q="$1"
  local d="$2"
  local ans=""
  read -r -p "$q [$d]: " ans || true
  echo "${ans:-$d}"
}

need_cmd go

# Safety: directory must be empty
if find . -mindepth 1 -maxdepth 1 -print -quit | grep -q .; then
  echo "Error: directory is not empty."
  echo "Run this installer in an empty directory."
  exit 1
fi

echo "== Ikaitla CLI Project Installer =="

MODULE_NAME="$(prompt_default \
  "Go module name (e.g. github.com/you/my-cli)" \
  "$(basename "$(pwd)")")"

FRAMEWORK_VERSION="$(prompt_default \
  "Ikaitla framework version" \
  "v0.1.3")"

echo ""
echo "-> Initializing Go module: $MODULE_NAME"
go mod init "$MODULE_NAME" >/dev/null

echo "-> Writing main.go"
cat > main.go <<'EOF'
package main

import (
	"fmt"

	"github.com/ikaitla/framework"
	_ "github.com/ikaitla/framework/all"
)

func main() {
	fmt.Println(framework.EngineName, framework.EngineVersion)
	if framework.EngineTagline != "" {
		fmt.Println(framework.EngineTagline)
	}
}
EOF

echo "-> Installing Ikaitla framework ($FRAMEWORK_VERSION)"
go get "github.com/ikaitla/framework@$FRAMEWORK_VERSION" >/dev/null

echo "-> Tidying modules"
go mod tidy

echo "-> Vendoring dependencies"
go mod vendor

echo "-> Building (vendor mode)"
go build -mod=vendor ./...

echo ""
echo "[GOOD] Ikaitla CLI project ready"
echo ""
echo "Next steps:"
echo "  - Add profiles under: cmd/<profile>/"
echo "  - Rebuild with: go build"
echo ""
echo "Update later with:"
echo "  go get github.com/ikaitla/framework@<new_tag>"
echo "  go mod tidy && go mod vendor"
