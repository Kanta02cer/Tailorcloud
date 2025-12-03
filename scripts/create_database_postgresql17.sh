#!/bin/bash

# PostgreSQL 17: tailorcloudデータベース作成スクリプト
# 使用方法: ./scripts/create_database_postgresql17.sh

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "PostgreSQL 17: tailorcloudデータベース作成"
echo "=========================================="
echo ""

# PostgreSQL 17のpsqlパス
PSQL="/Library/PostgreSQL/17/bin/psql"

if [ ! -f "$PSQL" ]; then
    echo -e "${RED}❌ PostgreSQL 17のpsqlが見つかりません: $PSQL${NC}"
    exit 1
fi

echo "PostgreSQL 17のpsql: $PSQL"
echo ""

# データベースを作成
echo -e "${BLUE}[1/2] データベース 'tailorcloud' を作成中...${NC}"
if $PSQL -d postgres -U postgres -c "CREATE DATABASE tailorcloud OWNER tailorcloud;" 2>&1; then
    echo -e "${GREEN}  ✅ データベースを作成しました${NC}"
elif echo "$?" | grep -q "already exists"; then
    echo -e "${YELLOW}  ⚠️  データベースは既に存在します${NC}"
else
    echo -e "${RED}  ❌ データベースの作成に失敗しました${NC}"
    exit 1
fi
echo ""

# 権限を付与
echo -e "${BLUE}[2/2] 権限を付与中...${NC}"
$PSQL -d tailorcloud -U postgres << 'SQL'
GRANT ALL PRIVILEGES ON DATABASE tailorcloud TO tailorcloud;
GRANT ALL ON SCHEMA public TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO tailorcloud;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO tailorcloud;
SQL

if [ $? -eq 0 ]; then
    echo -e "${GREEN}  ✅ 権限を付与しました${NC}"
else
    echo -e "${YELLOW}  ⚠️  権限の付与でエラーが発生しました（既に付与されている可能性があります）${NC}"
fi
echo ""

# 確認
echo "=========================================="
echo "確認中..."
echo "=========================================="
$PSQL -d tailorcloud -U postgres -c "SELECT current_database() AS database, current_user AS user;" 2>&1

echo ""
echo "=========================================="
echo -e "${GREEN}✅ セットアップ完了${NC}"
echo "=========================================="
echo ""
echo "次のステップ:"
echo "  1. 環境変数を設定: export POSTGRES_PASSWORD='tailorcloud_dev_password'"
echo "  2. マイグレーション実行: ./scripts/run_migrations_suit_mbti.sh"
echo ""

