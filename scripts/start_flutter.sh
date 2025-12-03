#!/bin/bash

# TailorCloud: Flutterアプリ起動スクリプト

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== TailorCloud Flutterアプリ起動 ==="
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 環境変数ファイルがあれば読み込む
ENV_FILE="$(pwd)/.env.local"
if [ -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}📝 環境変数ファイルを読み込み中: $ENV_FILE${NC}"
    export $(cat "$ENV_FILE" | grep -v '^#' | xargs)
fi

# デフォルト値の設定
export API_BASE_URL=${API_BASE_URL:-http://localhost:8080}

echo -e "${GREEN}✅ 環境変数:${NC}"
echo "  API_BASE_URL: $API_BASE_URL"
echo ""

# Flutterアプリディレクトリに移動
cd tailor-cloud-app

# Flutterがインストールされているか確認
if ! command -v flutter &> /dev/null; then
    echo -e "${RED}❌ Flutterがインストールされていません${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Flutter バージョン: $(flutter --version | head -n 1)${NC}"
echo ""

# 依存関係の確認
echo "📦 依存関係を確認中..."
flutter pub get
echo ""

# デバイスの確認
echo "📱 利用可能なデバイス:"
flutter devices
echo ""

# Flutterアプリを起動
echo -e "${GREEN}🚀 Flutterアプリを起動中...${NC}"
echo "  環境変数: API_BASE_URL=$API_BASE_URL"
echo ""

# 環境変数を指定してFlutterアプリを起動
flutter run --dart-define=API_BASE_URL="$API_BASE_URL"

