#!/bin/bash

# TailorCloud: バックエンドAPI起動スクリプト

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== TailorCloud バックエンドAPI起動 ==="
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 環境変数ファイルがあれば読み込む
ENV_FILE="$(pwd)/.env.local"
if [ -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}📝 環境変数ファイルを読み込み中: $ENV_FILE${NC}"
    export $(cat "$ENV_FILE" | grep -v '^#' | xargs)
fi

# デフォルト値の設定（既に設定されている場合は上書きしない）
export PORT=${PORT:-8080}
export POSTGRES_HOST=${POSTGRES_HOST:-localhost}
export POSTGRES_PORT=${POSTGRES_PORT:-5432}
export POSTGRES_USER=${POSTGRES_USER:-tailorcloud}
# POSTGRES_PASSWORDは.env.localから読み込まれるため、ここでは設定しない
# export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}
export POSTGRES_DB=${POSTGRES_DB:-tailorcloud}
export POSTGRES_SSLMODE=${POSTGRES_SSLMODE:-disable}

echo -e "${GREEN}✅ 環境変数:${NC}"
echo "  PORT: $PORT"
echo "  POSTGRES_HOST: $POSTGRES_HOST"
echo "  POSTGRES_PORT: $POSTGRES_PORT"
echo "  POSTGRES_USER: $POSTGRES_USER"
echo "  POSTGRES_DB: $POSTGRES_DB"
echo ""

# バックエンドディレクトリに移動
cd tailor-cloud-backend

# Goがインストールされているか確認
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Goがインストールされていません${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go バージョン: $(go version)${NC}"
echo ""

# 依存関係のインストール
echo "📦 依存関係をインストール中..."
go mod download
echo ""

# バックエンドAPIを起動
echo -e "${GREEN}🚀 バックエンドAPIを起動中...${NC}"
echo "  URL: http://localhost:$PORT"
echo "  Health Check: http://localhost:$PORT/health"
echo ""

go run cmd/api/main.go

