#!/bin/bash

# クイックテストスクリプト
# バックエンドとFlutterアプリの起動確認を簡単に行います

set -e

echo "🚀 TailorCloud クイックテスト"
echo ""

# バックエンドの起動確認
echo "📦 バックエンドの確認..."
cd "$(dirname "$0")/.."

# 環境変数の設定
export GCP_PROJECT_ID="${GCP_PROJECT_ID:-regalis-erp}"
export DEFAULT_TENANT_ID="${DEFAULT_TENANT_ID:-00000000-0000-0000-0000-000000000001}"

echo "   GCP_PROJECT_ID: $GCP_PROJECT_ID"
echo "   DEFAULT_TENANT_ID: $DEFAULT_TENANT_ID"
echo ""

# ビルド確認
echo "🔧 ビルド確認中..."
if go build ./cmd/api 2>/dev/null; then
    echo "   ✅ ビルド成功"
else
    echo "   ❌ ビルドエラー"
    exit 1
fi

# ポート確認
echo "🔌 ポート8080の確認..."
if lsof -i :8080 >/dev/null 2>&1; then
    echo "   ⚠️  ポート8080は既に使用されています"
    echo "   既存のプロセスを停止するか、別のポートを使用してください"
else
    echo "   ✅ ポート8080は利用可能"
fi

echo ""
echo "✅ バックエンドの準備が完了しました"
echo ""
echo "📝 次のコマンドでバックエンドを起動できます:"
echo "   ./scripts/start_backend_dev.sh"
echo ""
echo "または手動で:"
echo "   export GCP_PROJECT_ID=\"regalis-erp\""
echo "   export DEFAULT_TENANT_ID=\"00000000-0000-0000-0000-000000000001\""
echo "   go run cmd/api/main.go"

