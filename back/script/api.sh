#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
用法：
  ./script/api.sh          生成 API 代码并执行本地检查
  ./script/api.sh generate 只检查并生成 API 代码
  ./script/api.sh check    生成代码、测试、vet 并编译
  ./script/api.sh e2e      执行需要已启动服务的真实 HTTP 测试
EOF
}

generate() {
  echo "==> 检查 Proto"
  buf lint \
    --path api/common/v1 \
    --path api/business/v1 \
    --path api/order/v1

  echo "==> 生成 Go API 代码"
  buf generate \
    --path api/common/v1 \
    --path api/business/v1 \
    --path api/order/v1
}

unit_test() {
  echo "==> 执行不依赖已启动服务的测试"
  go test ./internal/... ./api/... ./cmd/...
}

check() {
  generate
  unit_test

  echo "==> 静态检查"
  go vet ./internal/... ./api/... ./cmd/...

  echo "==> 编译服务"
  go build ./cmd/business ./cmd/order
}

case "${1:-check}" in
  generate)
    generate
    ;;
  check)
    check
    ;;
  e2e)
    echo "==> 执行真实 HTTP 测试（请先启动 business 和 order 服务）"
    go test -v ./test/...
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
