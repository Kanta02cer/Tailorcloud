#!/bin/bash

# 開発環境用バックエンド起動スクリプト
# PostgreSQLがなくてもFirebase認証部分は動作します

set -e

echo "🚀 TailorCloud Backend 開発環境起動"
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 環境変数の設定（development.envから読み込む場合）
if [ -f "../tailor-cloud-app/config/development.env" ]; then
    echo "📝 development.envから環境変数を読み込み中..."
    source <(grep -v '^#' ../tailor-cloud-app/config/development.env | sed 's/^/export /')
fi

# デフォルト値の設定
export DEFAULT_TENANT_ID="${DEFAULT_TENANT_ID:-00000000-0000-0000-0000-000000000001}"
export GCP_PROJECT_ID="${GCP_PROJECT_ID:-regalis-erp}"
export PORT="${PORT:-8080}"

echo "📊 環境変数:"
echo "  DEFAULT_TENANT_ID: $DEFAULT_TENANT_ID"
echo "  GCP_PROJECT_ID: $GCP_PROJECT_ID"
echo "  PORT: $PORT"
echo ""

# Firebase認証情報の確認
if [ -z "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    echo "⚠️  GOOGLE_APPLICATION_CREDENTIALS が設定されていません"
    echo "   Firebase認証機能は制限される可能性があります"
    echo ""
fi

echo "🔧 バックエンドを起動中..."
echo "   URL: http://localhost:$PORT"
echo "   Health Check: http://localhost:$PORT/health"
echo ""

# バックエンドを起動
go run cmd/api/main.go

