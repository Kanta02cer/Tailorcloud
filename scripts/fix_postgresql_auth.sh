#!/bin/bash

# PostgreSQL認証設定修正スクリプト
# 使用方法: ./scripts/fix_postgresql_auth.sh

set +e  # エラー時に即座に終了しない

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "PostgreSQL認証設定修正"
echo "=========================================="
echo ""

# PostgreSQL PATH設定
source ./scripts/fix_postgresql_path.sh > /dev/null 2>&1

# PostgreSQLのデータディレクトリを確認
PG_DATA_DIR="/opt/homebrew/var/postgresql@15"
if [ ! -d "$PG_DATA_DIR" ]; then
    PG_DATA_DIR="/usr/local/var/postgresql@15"
fi

if [ ! -d "$PG_DATA_DIR" ]; then
    echo -e "${RED}❌ PostgreSQLデータディレクトリが見つかりません${NC}"
    exit 1
fi

PG_HBA_CONF="$PG_DATA_DIR/pg_hba.conf"

echo "PostgreSQLデータディレクトリ: $PG_DATA_DIR"
echo "認証設定ファイル: $PG_HBA_CONF"
echo ""

# pg_hba.confファイルが存在するか確認
if [ ! -f "$PG_HBA_CONF" ]; then
    echo -e "${RED}❌ pg_hba.confファイルが見つかりません${NC}"
    echo "PostgreSQLが正しくインストールされているか確認してください。"
    exit 1
fi

# バックアップを作成
if [ ! -f "$PG_HBA_CONF.backup" ]; then
    cp "$PG_HBA_CONF" "$PG_HBA_CONF.backup"
    echo -e "${GREEN}✅ バックアップを作成しました: $PG_HBA_CONF.backup${NC}"
else
    echo -e "${YELLOW}⚠️  バックアップは既に存在します${NC}"
fi

echo ""
echo "現在の認証設定を確認します..."
echo ""

# 現在の設定を表示
grep -E "^(local|host)" "$PG_HBA_CONF" | head -5 || echo "設定が見つかりません"

echo ""
echo "認証方式を 'trust' に変更しますか？ (y/n)"
echo "（これにより、パスワードなしで接続できるようになります）"
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo ""
    echo "認証設定を変更しています..."
    
    # 既存の設定をコメントアウト（安全のため）
    sed -i.bak -E 's/^(local|host)/# \1/' "$PG_HBA_CONF"
    
    # 新しい設定を追加
    cat >> "$PG_HBA_CONF" << 'AUTH_EOF'

# Trust authentication for local connections (added by fix_postgresql_auth.sh)
local   all             all                                     trust
host    all             all             127.0.0.1/32            trust
host    all             all             ::1/128                 trust
AUTH_EOF
    
    echo -e "${GREEN}✅ 認証設定を変更しました${NC}"
    echo ""
    echo "PostgreSQLを再起動します..."
    
    # PostgreSQLサービスを再起動
    brew services restart postgresql@15
    
    sleep 2
    
    echo ""
    echo "接続テストを行います..."
    
    # 接続テスト
    if psql -d postgres -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 接続テスト成功！${NC}"
        echo ""
        echo "次のステップ:"
        echo "  1. データベースを作成: createdb tailorcloud"
        echo "  2. ユーザーを作成: ./scripts/setup_postgresql_user_db.sh"
    else
        echo -e "${YELLOW}⚠️  接続テスト失敗（少し待ってから再試行してください）${NC}"
    fi
else
    echo "設定を変更せずに終了します。"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}✅ 完了${NC}"
echo "=========================================="
echo ""

set -e

