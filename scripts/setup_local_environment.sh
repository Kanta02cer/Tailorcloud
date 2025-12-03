#!/bin/bash

# TailorCloud: ローカル開発環境セットアップスクリプト
# このスクリプトは、ローカル環境でシステムを起動するための環境変数を設定します

set -e

echo "=== TailorCloud ローカル開発環境セットアップ ==="
echo ""

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# .env ファイルのパス
ENV_FILE="$(pwd)/.env.local"

echo "📝 環境変数ファイルを作成中: $ENV_FILE"
echo ""

# バックエンド環境変数
cat > "$ENV_FILE" << 'EOF'
# TailorCloud ローカル開発環境変数

# バックエンドAPI設定
PORT=8080

# PostgreSQL設定（オプション: PostgreSQLが不要な場合は空欄でも可）
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=tailorcloud
POSTGRES_PASSWORD=
POSTGRES_DB=tailorcloud
POSTGRES_SSLMODE=disable

# Firebase設定（オプション: デモ用には不要）
GCP_PROJECT_ID=
GOOGLE_APPLICATION_CREDENTIALS=

# Cloud Storage設定（オプション: デモ用には不要）
GCS_BUCKET_NAME=

# Flutterアプリ設定
API_BASE_URL=http://localhost:8080
EOF

echo -e "${GREEN}✅ 環境変数ファイルを作成しました: $ENV_FILE${NC}"
echo ""
echo "📋 設定内容:"
echo "  - バックエンドAPI: http://localhost:8080"
echo "  - PostgreSQL: オプション（未設定でも起動可能）"
echo "  - Firebase: オプション（未設定でも起動可能）"
echo ""
echo "🔧 環境変数を読み込むには、以下を実行してください:"
echo "  source $ENV_FILE"
echo "  または"
echo "  export \$(cat $ENV_FILE | grep -v '^#' | xargs)"
echo ""

