#!/bin/bash

# TailorCloud 本番環境ビルドスクリプト
# 使用方法: ./scripts/build_production.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG_DIR="$APP_DIR/config"

echo "🚀 TailorCloud 本番環境ビルドを開始します..."

# 設定ファイルの確認
if [ ! -f "$CONFIG_DIR/production.env" ]; then
    echo "⚠️  警告: $CONFIG_DIR/production.env が見つかりません"
    echo "📝 $CONFIG_DIR/production.env.example をコピーして設定してください"
    exit 1
fi

# 環境変数を読み込み
source "$CONFIG_DIR/production.env"

# 必須設定の確認
if [ -z "$API_BASE_URL" ]; then
    echo "❌ エラー: API_BASE_URL が設定されていません"
    exit 1
fi

echo "📋 ビルド設定:"
echo "   - 環境: $ENV"
echo "   - API URL: $API_BASE_URL"
echo "   - Firebase: $([ "$ENABLE_FIREBASE" = "true" ] && echo "有効" || echo "無効")"

# Flutterアプリディレクトリに移動
cd "$APP_DIR"

# 依存関係の取得
echo "📦 依存関係を取得しています..."
flutter pub get

# ビルド引数の構築
BUILD_ARGS=(
    "--dart-define=ENV=$ENV"
    "--dart-define=API_BASE_URL=$API_BASE_URL"
    "--dart-define=DEFAULT_TENANT_ID=${DEFAULT_TENANT_ID:-tenant-production-001}"
)

# Firebase設定がある場合は追加
if [ "$ENABLE_FIREBASE" = "true" ]; then
    if [ -n "$FIREBASE_API_KEY" ] && [ -n "$FIREBASE_APP_ID" ] && [ -n "$FIREBASE_PROJECT_ID" ]; then
        BUILD_ARGS+=(
            "--dart-define=ENABLE_FIREBASE=true"
            "--dart-define=FIREBASE_API_KEY=$FIREBASE_API_KEY"
            "--dart-define=FIREBASE_APP_ID=$FIREBASE_APP_ID"
            "--dart-define=FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID"
        )
        if [ -n "$FIREBASE_MESSAGING_SENDER_ID" ]; then
            BUILD_ARGS+=("--dart-define=FIREBASE_MESSAGING_SENDER_ID=$FIREBASE_MESSAGING_SENDER_ID")
        fi
    else
        echo "⚠️  警告: Firebaseが有効ですが、設定が不完全です。Firebaseなしでビルドします。"
    fi
fi

# Webビルド
echo "🌐 Webアプリをビルドしています..."
flutter build web --release "${BUILD_ARGS[@]}"

echo "✅ ビルドが完了しました！"
echo "📁 ビルド成果物: $APP_DIR/build/web"

