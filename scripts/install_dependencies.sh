#!/bin/bash

# TailorCloud: 依存関係インストールスクリプト
# このスクリプトは、システムに必要なすべての依存関係をインストールします

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== TailorCloud 依存関係インストール ==="
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 1. バックエンド依存関係のインストール
echo -e "${GREEN}📦 バックエンド依存関係をインストール中...${NC}"
cd tailor-cloud-backend

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Goがインストールされていません${NC}"
    echo "  インストール方法: brew install go"
    exit 1
fi

go mod download
go mod tidy
echo -e "${GREEN}✅ バックエンド依存関係のインストール完了${NC}"
echo ""

# 2. Flutterアプリ依存関係のインストール
echo -e "${GREEN}📦 Flutterアプリ依存関係をインストール中...${NC}"
cd ../tailor-cloud-app

if ! command -v flutter &> /dev/null; then
    echo -e "${RED}❌ Flutterがインストールされていません${NC}"
    echo "  インストール方法: brew install --cask flutter"
    exit 1
fi

flutter pub get
echo -e "${GREEN}✅ Flutterアプリ依存関係のインストール完了${NC}"
echo ""

# 3. Flutterコード生成
echo -e "${GREEN}🔧 Flutterコード生成を実行中...${NC}"
flutter pub run build_runner build --delete-conflicting-outputs
echo -e "${GREEN}✅ Flutterコード生成完了${NC}"
echo ""

# 4. Webアプリ依存関係のインストール（オプション）
echo -e "${GREEN}📦 Webアプリ依存関係をインストール中...${NC}"
cd ../suit-mbti-web-app

if ! command -v npm &> /dev/null; then
    echo -e "${YELLOW}⚠️  npmがインストールされていません（オプション）${NC}"
    echo "  インストール方法: brew install node"
else
    npm install
    echo -e "${GREEN}✅ Webアプリ依存関係のインストール完了${NC}"
fi
echo ""

# プロジェクトルートに戻る
cd ..

echo -e "${GREEN}=== 依存関係インストール完了 ===${NC}"
echo ""
echo "次のステップ:"
echo "  1. 環境変数のセットアップ: ./scripts/setup_local_environment.sh"
echo "  2. システム状態の確認: ./scripts/check_system.sh"
echo "  3. バックエンドの起動: ./scripts/start_backend.sh"
echo "  4. Flutterアプリの起動: ./scripts/start_flutter.sh"
echo ""

