#!/bin/bash

# Suit-MBTI統合用マイグレーション実行スクリプト
# 使用方法: ./scripts/run_migrations_suit_mbti.sh

set +e  # エラー時に即座に終了しない（接続確認のため）

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Suit-MBTI統合用マイグレーション実行"
echo "=========================================="
echo ""

# 環境変数の読み込み
if [ -f .env.local ]; then
    source .env.local
    echo "✅ .env.localファイルを読み込みました"
elif [ -f .env ]; then
    source .env
    echo "✅ .envファイルを読み込みました"
else
    echo -e "${YELLOW}⚠️  環境変数ファイル（.env.local または .env）が見つかりません。${NC}"
    echo -e "${YELLOW}   環境変数を直接設定するか、scripts/setup_local_environment.sh を実行してください。${NC}"
fi

# PostgreSQL接続情報の確認
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
POSTGRES_PORT=${POSTGRES_PORT:-5432}
POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}
POSTGRES_DB=${POSTGRES_DB:-tailorcloud}

echo "データベース接続情報:"
echo "  Host: $POSTGRES_HOST"
echo "  Port: $POSTGRES_PORT"
echo "  User: $POSTGRES_USER"
echo "  Database: $POSTGRES_DB"
echo ""

# PostgreSQL接続確認
echo "PostgreSQL接続を確認しています..."
export PGPASSWORD=$POSTGRES_PASSWORD

# パスワードが空の場合は、パスワードなしで接続を試みる
if [ -z "$POSTGRES_PASSWORD" ]; then
    echo -e "${YELLOW}⚠️  POSTGRES_PASSWORDが設定されていません。パスワードなしで接続を試みます。${NC}"
fi

if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PostgreSQL接続成功${NC}"
else
    echo -e "${RED}❌ PostgreSQL接続失敗${NC}"
    echo ""
    echo "接続情報を確認してください:"
    echo "  Host: $POSTGRES_HOST"
    echo "  Port: $POSTGRES_PORT"
    echo "  User: $POSTGRES_USER"
    echo "  Database: $POSTGRES_DB"
    echo ""
    echo "トラブルシューティング:"
    echo "  1. PostgreSQLが起動しているか確認: pg_isready"
    echo "  2. データベースが存在するか確認: psql -l"
    echo "  3. ユーザーが存在するか確認: psql -U $POSTGRES_USER -d postgres -c '\du'"
    echo "  4. パスワードを設定: export POSTGRES_PASSWORD='your_password'"
    echo ""
    echo "環境変数を直接設定して再実行する例:"
    echo "  export POSTGRES_HOST=localhost"
    echo "  export POSTGRES_PORT=5432"
    echo "  export POSTGRES_USER=your_user"
    echo "  export POSTGRES_PASSWORD=your_password"
    echo "  export POSTGRES_DB=tailorcloud"
    echo "  ./scripts/run_migrations_suit_mbti.sh"
    echo ""
    exit 1
fi
echo ""

# エラー時に即座に終了するように戻す
set -e

# マイグレーションファイルのパス
MIGRATIONS_DIR="tailor-cloud-backend/migrations"

# マイグレーションファイルのリスト
MIGRATIONS=(
    "013_create_diagnoses_table.sql"
    "014_create_appointments_table.sql"
    "015_extend_customers_for_suit_mbti.sql"
)

echo "マイグレーションを実行します..."
echo ""

# 各マイグレーションを実行
for migration in "${MIGRATIONS[@]}"; do
    MIGRATION_PATH="$MIGRATIONS_DIR/$migration"
    
    if [ ! -f "$MIGRATION_PATH" ]; then
        echo -e "${RED}❌ マイグレーションファイルが見つかりません: $MIGRATION_PATH${NC}"
        exit 1
    fi
    
    echo "実行中: $migration"
    
    if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f "$MIGRATION_PATH" > /dev/null 2>&1; then
        echo -e "${GREEN}  ✅ 成功${NC}"
    else
        # エラーが発生した場合でも、既に存在する場合はスキップ
        ERROR_OUTPUT=$(psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f "$MIGRATION_PATH" 2>&1)
        if echo "$ERROR_OUTPUT" | grep -q "already exists"; then
            echo -e "${YELLOW}  ⚠️  テーブルは既に存在します（スキップ）${NC}"
        else
            echo -e "${RED}  ❌ 失敗${NC}"
            echo "$ERROR_OUTPUT"
            exit 1
        fi
    fi
    echo ""
done

echo "=========================================="
echo -e "${GREEN}✅ マイグレーション完了${NC}"
echo "=========================================="
echo ""
echo "作成されたテーブル:"
echo "  - diagnoses (診断ログ)"
echo "  - appointments (予約管理)"
echo "  - customers (拡張: occupation, ltv_score等)"
echo ""
echo "次のステップ:"
echo "  1. バックエンドサーバーを起動"
echo "  2. API動作テストを実行"
echo "  3. ドキュメント参照: docs/78_Suit_MBTI_Feature_Guide.md"
echo ""

