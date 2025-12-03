#!/bin/bash

# PostgreSQL PATH修正スクリプト
# 使用方法: source ./scripts/fix_postgresql_path.sh

echo "=== PostgreSQL PATH設定 ==="
echo ""

# HomebrewのPostgreSQLのパスを確認
PG_BREW_PREFIX=$(brew --prefix postgresql@15 2>/dev/null)

if [ -z "$PG_BREW_PREFIX" ]; then
    echo "⚠️  PostgreSQL@15が見つかりません。"
    echo ""
    echo "インストール方法:"
    echo "  brew install postgresql@15"
    exit 1
fi

PG_BIN_PATH="$PG_BREW_PREFIX/bin"

if [ -d "$PG_BIN_PATH" ]; then
    echo "✅ PostgreSQLのパスを確認: $PG_BIN_PATH"
    
    # PATHに追加（まだ含まれていない場合）
    if [[ ":$PATH:" != *":$PG_BIN_PATH:"* ]]; then
        export PATH="$PG_BIN_PATH:$PATH"
        echo "✅ PATHに追加しました: $PG_BIN_PATH"
    else
        echo "✅ PATHに既に含まれています"
    fi
    
    echo ""
    echo "psqlの場所:"
    which psql
    
    echo ""
    echo "psqlバージョン:"
    psql --version
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✅ 設定完了！"
    echo ""
    echo "この設定を永続化するには、以下を~/.zshrcに追加してください:"
    echo "  export PATH=\"$PG_BIN_PATH:\$PATH\""
    echo ""
    echo "または、このスクリプトをsourceして使ってください:"
    echo "  source ./scripts/fix_postgresql_path.sh"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
else
    echo "❌ PostgreSQLのbinディレクトリが見つかりません: $PG_BIN_PATH"
    exit 1
fi

