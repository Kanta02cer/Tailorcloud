#!/bin/bash

# TailorCloud 開発環境実行スクリプト
# 使用方法: ./scripts/run_development.sh [device]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG_DIR="$APP_DIR/config"

DEVICE="${1:-chrome}"

echo "🚀 TailorCloud 開発環境を起動します..."

# 設定ファイルの確認（オプショナル）
if [ -f "$CONFIG_DIR/development.env" ]; then
    echo "📝 設定ファイルを読み込みます..."
    source "$CONFIG_DIR/development.env"
else
    echo "📝 デフォルト設定を使用します..."
    ENV="development"
    API_BASE_URL="http://localhost:8080"
    ENABLE_FIREBASE="false"
    DEFAULT_TENANT_ID="tenant-123"
fi

echo "📋 実行設定:"
echo "   - 環境: ${ENV:-development}"
echo "   - API URL: ${API_BASE_URL:-http://localhost:8080}"
echo "   - デバイス: $DEVICE"
echo "   - Firebase: $([ "${ENABLE_FIREBASE:-false}" = "true" ] && echo "有効" || echo "無効")"

# Flutterアプリディレクトリに移動
cd "$APP_DIR"

# 依存関係の取得
echo "📦 依存関係を取得しています..."
flutter pub get

# 実行引数の構築
RUN_ARGS=(
    "--dart-define=ENV=${ENV:-development}"
    "--dart-define=API_BASE_URL=${API_BASE_URL:-http://localhost:8080}"
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
        # オプショナルなFirebase設定
        if [ -n "$FIREBASE_AUTH_DOMAIN" ]; then
            RUN_ARGS+=("--dart-define=FIREBASE_AUTH_DOMAIN=$FIREBASE_AUTH_DOMAIN")
        fi
        if [ -n "$FIREBASE_STORAGE_BUCKET" ]; then
            RUN_ARGS+=("--dart-define=FIREBASE_STORAGE_BUCKET=$FIREBASE_STORAGE_BUCKET")
        fi
        if [ -n "$FIREBASE_MEASUREMENT_ID" ]; then
            RUN_ARGS+=("--dart-define=FIREBASE_MEASUREMENT_ID=$FIREBASE_MEASUREMENT_ID")
        fi
    fi
fi

# アプリを実行
echo "🎯 アプリを起動しています..."
flutter run -d "$DEVICE" "${RUN_ARGS[@]}"

