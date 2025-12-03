#!/bin/bash

# TailorCloud: システム状態確認スクリプト

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== TailorCloud システム状態確認 ==="
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 1. 必要なコマンドの確認
echo "🔍 必要なコマンドの確認:"
echo ""

COMMANDS=("go" "flutter" "docker" "psql")
for cmd in "${COMMANDS[@]}"; do
    if command -v $cmd &> /dev/null; then
        echo -e "  ${GREEN}✅ $cmd: インストール済み${NC}"
        case $cmd in
            go)
                echo "     バージョン: $(go version | awk '{print $3}')"
                ;;
            flutter)
                echo "     バージョン: $(flutter --version | head -n 1 | awk '{print $2}')"
                ;;
        esac
    else
        echo -e "  ${RED}❌ $cmd: インストールされていません${NC}"
    fi
done
echo ""

# 2. 環境変数ファイルの確認
echo "📝 環境変数ファイルの確認:"
ENV_FILE="$(pwd)/.env.local"
if [ -f "$ENV_FILE" ]; then
    echo -e "  ${GREEN}✅ 環境変数ファイルが見つかりました: $ENV_FILE${NC}"
else
    echo -e "  ${YELLOW}⚠️  環境変数ファイルが見つかりません: $ENV_FILE${NC}"
    echo "     セットアップスクリプトを実行してください: ./scripts/setup_local_environment.sh"
fi
echo ""

# 3. バックエンドの確認
echo "🔧 バックエンドの確認:"
if [ -d "tailor-cloud-backend" ]; then
    echo -e "  ${GREEN}✅ バックエンドディレクトリが見つかりました${NC}"
    cd tailor-cloud-backend
    
    # Go依存関係の確認
    if [ -f "go.mod" ]; then
        echo -e "  ${GREEN}✅ go.mod が見つかりました${NC}"
    else
        echo -e "  ${RED}❌ go.mod が見つかりません${NC}"
    fi
    
    cd ..
else
    echo -e "  ${RED}❌ バックエンドディレクトリが見つかりません${NC}"
fi
echo ""

# 4. Flutterアプリの確認
echo "📱 Flutterアプリの確認:"
if [ -d "tailor-cloud-app" ]; then
    echo -e "  ${GREEN}✅ Flutterアプリディレクトリが見つかりました${NC}"
    cd tailor-cloud-app
    
    # pubspec.yamlの確認
    if [ -f "pubspec.yaml" ]; then
        echo -e "  ${GREEN}✅ pubspec.yaml が見つかりました${NC}"
    else
        echo -e "  ${RED}❌ pubspec.yaml が見つかりません${NC}"
    fi
    
    cd ..
else
    echo -e "  ${RED}❌ Flutterアプリディレクトリが見つかりません${NC}"
fi
echo ""

# 5. PostgreSQLの確認（オプション）
echo "🗄️  PostgreSQLの確認（オプション）:"
if command -v psql &> /dev/null; then
    echo -e "  ${GREEN}✅ psql コマンドが利用可能です${NC}"
    echo "     接続テストを実行するには: psql -h localhost -U tailorcloud -d tailorcloud"
else
    echo -e "  ${YELLOW}⚠️  PostgreSQLはオプションです（Firestoreモードで動作可能）${NC}"
fi
echo ""

# 6. バックエンドAPIのヘルスチェック
echo "🌐 バックエンドAPIのヘルスチェック:"
if command -v curl &> /dev/null; then
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅ バックエンドAPIが起動しています (http://localhost:8080)${NC}"
    else
        echo -e "  ${YELLOW}⚠️  バックエンドAPIが起動していません${NC}"
        echo "     起動するには: ./scripts/start_backend.sh"
    fi
else
    echo -e "  ${YELLOW}⚠️  curl コマンドが利用できません（手動確認が必要）${NC}"
fi
echo ""

echo "=== システム状態確認完了 ==="

