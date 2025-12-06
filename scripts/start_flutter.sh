#!/bin/bash

# TailorCloud: Flutterアプリ起動スクリプト（改善版）
# 使用方法: ./scripts/start_flutter.sh [environment] [device]
# 例: ./scripts/start_flutter.sh development chrome
# 例: ./scripts/start_flutter.sh production web-server

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 引数の取得
ENVIRONMENT="${1:-development}"
DEVICE="${2:-chrome}"

echo -e "${BLUE}=== TailorCloud Flutterアプリ起動 ===${NC}"
echo ""

# プロジェクトルートに移動
cd "$(dirname "$0")/.."

# 設定ファイルのパス
CONFIG_DIR="tailor-cloud-app/config"
ENV_FILE="$CONFIG_DIR/${ENVIRONMENT}.env"

# 環境変数ファイルの確認
if [ -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}📝 環境変数ファイルを読み込み中: $ENV_FILE${NC}"
    set -a
    source "$ENV_FILE"
    set +a
else
    echo -e "${YELLOW}⚠️  環境変数ファイルが見つかりません: $ENV_FILE${NC}"
    echo -e "${YELLOW}📝 デフォルト設定を使用します${NC}"
    
    # デフォルト値の設定
    export ENV="${ENVIRONMENT}"
    export API_BASE_URL="http://localhost:8080"
    export ENABLE_FIREBASE="false"
    export DEFAULT_TENANT_ID="tenant-123"
fi

# 必須設定の確認
if [ -z "$API_BASE_URL" ]; then
    echo -e "${RED}❌ エラー: API_BASE_URL が設定されていません${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 環境設定:${NC}"
echo "  環境: ${ENV:-$ENVIRONMENT}"
echo "  API URL: $API_BASE_URL"
echo "  Firebase: $([ "${ENABLE_FIREBASE:-false}" = "true" ] && echo "有効" || echo "無効")"
echo "  テナントID: ${DEFAULT_TENANT_ID:-tenant-123}"
echo ""

# Flutterアプリディレクトリに移動
cd tailor-cloud-app

# Flutterがインストールされているか確認
if ! command -v flutter &> /dev/null; then
    echo -e "${RED}❌ Flutterがインストールされていません${NC}"
    echo "  インストール方法: https://flutter.dev/docs/get-started/install"
    exit 1
fi

echo -e "${GREEN}✅ Flutter バージョン: $(flutter --version | head -n 1)${NC}"
echo ""

# 依存関係の確認
echo "📦 依存関係を確認中..."
flutter pub get
echo ""

# デバイスの確認（実行時のみ）
if [ "$DEVICE" != "web-server" ]; then
    echo "📱 利用可能なデバイス:"
    flutter devices
    echo ""
fi

# 実行引数の構築
RUN_ARGS=(
    "--dart-define=ENV=${ENV:-$ENVIRONMENT}"
    "--dart-define=API_BASE_URL=$API_BASE_URL"
    "--dart-define=DEFAULT_TENANT_ID=${DEFAULT_TENANT_ID:-tenant-123}"
)

# Firebase設定がある場合は追加
if [ "${ENABLE_FIREBASE:-false}" = "true" ]; then
    if [ -n "$FIREBASE_API_KEY" ] && [ -n "$FIREBASE_APP_ID" ] && [ -n "$FIREBASE_PROJECT_ID" ]; then
        RUN_ARGS+=(
            "--dart-define=ENABLE_FIREBASE=true"
            "--dart-define=FIREBASE_API_KEY=$FIREBASE_API_KEY"
            "--dart-define=FIREBASE_APP_ID=$FIREBASE_APP_ID"
            "--dart-define=FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID"
        )
        if [ -n "$FIREBASE_MESSAGING_SENDER_ID" ]; then
            RUN_ARGS+=("--dart-define=FIREBASE_MESSAGING_SENDER_ID=$FIREBASE_MESSAGING_SENDER_ID")
        fi
        echo -e "${GREEN}✅ Firebase設定を読み込みました${NC}"
    else
        echo -e "${YELLOW}⚠️  Firebaseが有効ですが、設定が不完全です。Firebaseなしで実行します。${NC}"
    fi
fi

# アプリを起動
echo -e "${GREEN}🚀 Flutterアプリを起動中...${NC}"
echo "  環境: ${ENV:-$ENVIRONMENT}"
echo "  デバイス: $DEVICE"
echo ""

if [ "$DEVICE" = "web-server" ]; then
    # Webサーバーとして起動（本番環境用）
    echo "🌐 Webサーバーを起動しています..."
    flutter run -d chrome --web-port=8080 --web-hostname=0.0.0.0 "${RUN_ARGS[@]}"
else
    # 通常の実行
    flutter run -d "$DEVICE" "${RUN_ARGS[@]}"
fi
