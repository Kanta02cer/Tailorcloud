#!/bin/bash

# PostgreSQLユーザーとデータベース作成スクリプト
# 使用方法: ./scripts/setup_postgresql_user_db.sh

set +e  # エラー時に即座に終了しない

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "PostgreSQLユーザー・データベース作成"
echo "=========================================="
echo ""

# PostgreSQL PATH設定
source ./scripts/fix_postgresql_path.sh > /dev/null 2>&1

# 環境変数の読み込み
if [ -f .env.local ]; then
    source .env.local
    echo "✅ .env.localファイルを読み込みました"
elif [ -f .env ]; then
    source .env
    echo "✅ .envファイルを読み込みました"
fi

# デフォルト値
POSTGRES_USER=${POSTGRES_USER:-tailorcloud}
POSTGRES_DB=${POSTGRES_DB:-tailorcloud}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}

echo "設定値:"
echo "  ユーザー名: $POSTGRES_USER"
echo "  データベース名: $POSTGRES_DB"
echo ""

# PostgreSQLに接続できるか確認
echo -e "${BLUE}[1/3] PostgreSQL接続確認...${NC}"
if psql -h localhost -U postgres -d postgres -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}  ✅ PostgreSQLに接続できました${NC}"
else
    echo -e "${RED}  ❌ PostgreSQLに接続できませんでした${NC}"
    echo ""
    echo "PostgreSQLが起動しているか確認してください:"
    echo "  brew services start postgresql@15"
    exit 1
fi
echo ""

# ユーザーの存在確認と作成
echo -e "${BLUE}[2/3] ユーザー '$POSTGRES_USER' の確認・作成...${NC}"
USER_EXISTS=$(psql -h localhost -U postgres -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$POSTGRES_USER'")
if [ "$USER_EXISTS" = "1" ]; then
    echo -e "${GREEN}  ✅ ユーザー '$POSTGRES_USER' は既に存在します${NC}"
else
    echo "  ユーザーを作成しています..."
    if [ -z "$POSTGRES_PASSWORD" ]; then
        psql -h localhost -U postgres -d postgres -c "CREATE USER $POSTGRES_USER;" 2>&1
    else
        psql -h localhost -U postgres -d postgres -c "CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';" 2>&1
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✅ ユーザーを作成しました${NC}"
    else
        echo -e "${YELLOW}  ⚠️  ユーザーの作成に失敗しました（既に存在する可能性があります）${NC}"
    fi
fi
echo ""

# データベースの存在確認と作成
echo -e "${BLUE}[3/3] データベース '$POSTGRES_DB' の確認・作成...${NC}"
DB_EXISTS=$(psql -h localhost -U postgres -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB'")
if [ "$DB_EXISTS" = "1" ]; then
    echo -e "${GREEN}  ✅ データベース '$POSTGRES_DB' は既に存在します${NC}"
else
    echo "  データベースを作成しています..."
    psql -h localhost -U postgres -d postgres -c "CREATE DATABASE $POSTGRES_DB OWNER $POSTGRES_USER;" 2>&1
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✅ データベースを作成しました${NC}"
    else
        echo -e "${RED}  ❌ データベースの作成に失敗しました${NC}"
        exit 1
    fi
fi

# 権限の付与
echo "  権限を付与しています..."
psql -h localhost -U postgres -d $POSTGRES_DB -c "GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;" > /dev/null 2>&1
psql -h localhost -U postgres -d $POSTGRES_DB -c "GRANT ALL ON SCHEMA public TO $POSTGRES_USER;" > /dev/null 2>&1
psql -h localhost -U postgres -d $POSTGRES_DB -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $POSTGRES_USER;" > /dev/null 2>&1
echo -e "${GREEN}  ✅ 権限を付与しました${NC}"
echo ""

# 接続テスト
echo "=========================================="
echo "接続テスト..."
echo "=========================================="
if [ -z "$POSTGRES_PASSWORD" ]; then
    if psql -h localhost -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 接続テスト成功！${NC}"
    else
        echo -e "${YELLOW}⚠️  接続テスト失敗（パスワードが必要な可能性があります）${NC}"
    fi
else
    export PGPASSWORD=$POSTGRES_PASSWORD
    if psql -h localhost -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 接続テスト成功！${NC}"
    else
        echo -e "${RED}❌ 接続テスト失敗${NC}"
        echo "パスワードを確認してください。"
    fi
fi

echo ""
echo "=========================================="
echo -e "${GREEN}✅ セットアップ完了${NC}"
echo "=========================================="
echo ""
echo "次のステップ:"
echo "  1. マイグレーション実行: ./scripts/run_migrations_suit_mbti.sh"
echo "  2. テストデータ準備: psql -f scripts/prepare_test_data_suit_mbti.sql"
echo ""

set -e

