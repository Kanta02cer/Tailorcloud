#!/bin/bash

# PostgreSQL接続確認スクリプト
# 使用方法: ./scripts/check_postgres_connection.sh

set +e  # エラー時に即座に終了しない

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "PostgreSQL接続確認"
echo "=========================================="
echo ""

# 環境変数の読み込み
if [ -f .env.local ]; then
    source .env.local
    echo -e "${GREEN}✅ .env.localファイルを読み込みました${NC}"
elif [ -f .env ]; then
    source .env
    echo -e "${GREEN}✅ .envファイルを読み込みました${NC}"
else
    echo -e "${YELLOW}⚠️  環境変数ファイルが見つかりません。${NC}"
    echo "デフォルト値を使用します。"
fi

# PostgreSQL接続情報の確認
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
POSTGRES_PORT=${POSTGRES_PORT:-5432}
POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}
POSTGRES_DB=${POSTGRES_DB:-tailorcloud}

echo "接続情報:"
echo "  Host: $POSTGRES_HOST"
echo "  Port: $POSTGRES_PORT"
echo "  User: $POSTGRES_USER"
echo "  Database: $POSTGRES_DB"
echo ""

# 1. PostgreSQLが起動しているか確認
echo -e "${BLUE}[1/4] PostgreSQLが起動しているか確認...${NC}"
if command -v pg_isready > /dev/null 2>&1; then
    if pg_isready -h $POSTGRES_HOST -p $POSTGRES_PORT > /dev/null 2>&1; then
        echo -e "${GREEN}  ✅ PostgreSQLは起動しています${NC}"
    else
        echo -e "${RED}  ❌ PostgreSQLが起動していません${NC}"
        echo ""
        echo "解決方法:"
        echo "  # macOS (Homebrew)"
        echo "  brew services start postgresql"
        echo ""
        echo "  # Linux"
        echo "  sudo systemctl start postgresql"
        echo ""
        exit 1
    fi
else
    echo -e "${YELLOW}  ⚠️  pg_isreadyコマンドが見つかりません。スキップします。${NC}"
fi
echo ""

# 2. データベース一覧を表示
echo -e "${BLUE}[2/4] データベース一覧を確認...${NC}"
export PGPASSWORD=$POSTGRES_PASSWORD
if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d postgres -c "\l" > /dev/null 2>&1; then
    echo -e "${GREEN}  ✅ データベース一覧を取得しました${NC}"
    psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d postgres -c "\l" | grep -E "Name|tailorcloud" || true
else
    echo -e "${RED}  ❌ データベース一覧を取得できませんでした${NC}"
    echo "接続情報を確認してください。"
    exit 1
fi
echo ""

# 3. 対象データベースが存在するか確認
echo -e "${BLUE}[3/4] データベース '$POSTGRES_DB' が存在するか確認...${NC}"
DB_EXISTS=$(psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB'")
if [ "$DB_EXISTS" = "1" ]; then
    echo -e "${GREEN}  ✅ データベース '$POSTGRES_DB' は存在します${NC}"
else
    echo -e "${YELLOW}  ⚠️  データベース '$POSTGRES_DB' が存在しません${NC}"
    echo ""
    echo "データベースを作成しますか？ (y/n)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        if createdb -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER $POSTGRES_DB 2>/dev/null; then
            echo -e "${GREEN}  ✅ データベースを作成しました${NC}"
        else
            echo -e "${RED}  ❌ データベースの作成に失敗しました${NC}"
            echo "手動で作成してください:"
            echo "  createdb -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER $POSTGRES_DB"
        fi
    else
        echo "データベースを作成せずに終了します。"
        exit 1
    fi
fi
echo ""

# 4. データベースに接続できるか確認
echo -e "${BLUE}[4/4] データベース '$POSTGRES_DB' に接続できるか確認...${NC}"
if psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}  ✅ データベースに接続できました${NC}"
    echo ""
    echo "=========================================="
    echo -e "${GREEN}✅ PostgreSQL接続確認完了${NC}"
    echo "=========================================="
    echo ""
    echo "次のステップ:"
    echo "  ./scripts/run_migrations_suit_mbti.sh"
    echo ""
else
    echo -e "${RED}  ❌ データベースに接続できませんでした${NC}"
    echo ""
    echo "トラブルシューティング:"
    echo "  1. ユーザーが存在するか確認: psql -U $POSTGRES_USER -d postgres -c '\\du'"
    echo "  2. ユーザーに権限があるか確認"
    echo "  3. パスワードが正しいか確認"
    echo ""
    echo "詳細は docs/83_Troubleshooting_Guide.md を参照してください。"
    exit 1
fi

