#!/bin/bash

# 認証システムのセットアップスクリプト
# デフォルトテナントの作成と環境変数の確認を行います

set -e

echo "🔐 TailorCloud 認証システム セットアップ"
echo ""

# データベース接続情報
DB_USER="${POSTGRES_USER:-tailorcloud}"
DB_NAME="${POSTGRES_DB:-tailorcloud}"
DB_HOST="${POSTGRES_HOST:-localhost}"
DB_PORT="${POSTGRES_PORT:-5432}"

echo "📊 データベース接続情報:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# データベース接続確認
echo "🔍 データベース接続を確認中..."
if ! PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ データベースに接続できません"
    echo "   接続情報を確認してください:"
    echo "   - POSTGRES_HOST"
    echo "   - POSTGRES_PORT"
    echo "   - POSTGRES_USER"
    echo "   - POSTGRES_PASSWORD"
    echo "   - POSTGRES_DB"
    exit 1
fi
echo "✅ データベース接続成功"
echo ""

# デフォルトテナントの作成
echo "🏢 デフォルトテナントを作成中..."
DEFAULT_TENANT_ID="00000000-0000-0000-0000-000000000001"

if PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f scripts/create_default_tenant.sql > /dev/null 2>&1; then
    echo "✅ デフォルトテナントを作成しました"
    echo "   Tenant ID: $DEFAULT_TENANT_ID"
else
    echo "⚠️  テナント作成に失敗しました（既に存在する可能性があります）"
fi
echo ""

# テナントの確認
echo "🔍 テナント情報を確認中..."
TENANT_INFO=$(PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT name || ' (' || type || ')' FROM tenants WHERE id = '$DEFAULT_TENANT_ID';" 2>/dev/null | xargs)

if [ -n "$TENANT_INFO" ]; then
    echo "✅ テナントが見つかりました: $TENANT_INFO"
else
    echo "❌ テナントが見つかりませんでした"
    exit 1
fi
echo ""

# 環境変数の確認
echo "📝 環境変数の確認:"
echo ""

if [ -z "$DEFAULT_TENANT_ID_ENV" ]; then
    echo "⚠️  DEFAULT_TENANT_ID が設定されていません"
    echo "   以下のコマンドで設定してください:"
    echo "   export DEFAULT_TENANT_ID=\"$DEFAULT_TENANT_ID\""
else
    echo "✅ DEFAULT_TENANT_ID: $DEFAULT_TENANT_ID_ENV"
fi

if [ -z "$GCP_PROJECT_ID" ]; then
    echo "⚠️  GCP_PROJECT_ID が設定されていません"
else
    echo "✅ GCP_PROJECT_ID: $GCP_PROJECT_ID"
fi

if [ -z "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    echo "⚠️  GOOGLE_APPLICATION_CREDENTIALS が設定されていません"
else
    echo "✅ GOOGLE_APPLICATION_CREDENTIALS: $GOOGLE_APPLICATION_CREDENTIALS"
fi
echo ""

# セットアップ完了
echo "✨ セットアップ完了！"
echo ""
echo "次のステップ:"
echo "1. 環境変数を設定:"
echo "   export DEFAULT_TENANT_ID=\"$DEFAULT_TENANT_ID\""
echo "   export GCP_PROJECT_ID=\"your-firebase-project-id\""
echo ""
echo "2. バックエンドを起動:"
echo "   go run cmd/api/main.go"
echo ""
echo "3. 動作確認:"
echo "   curl http://localhost:8080/health"
echo ""

