#!/bin/bash

# TailorCloud 環境設定セットアップスクリプト
# 使用方法: ./scripts/setup_environment.sh [environment]
# 環境: development, staging, production

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
CONFIG_DIR="$APP_DIR/config"

ENVIRONMENT="${1:-development}"

echo "🔧 TailorCloud 環境設定をセットアップします..."
echo "   環境: $ENVIRONMENT"

# 設定ディレクトリの作成
mkdir -p "$CONFIG_DIR"

# 環境に応じた設定ファイルのコピー
case "$ENVIRONMENT" in
    development)
        if [ ! -f "$CONFIG_DIR/development.env" ]; then
            if [ -f "$CONFIG_DIR/development.env.example" ]; then
                cp "$CONFIG_DIR/development.env.example" "$CONFIG_DIR/development.env"
                echo "✅ $CONFIG_DIR/development.env を作成しました"
            else
                echo "⚠️  $CONFIG_DIR/development.env.example が見つかりません"
            fi
        else
            echo "ℹ️  $CONFIG_DIR/development.env は既に存在します"
        fi
        ;;
    staging)
        if [ ! -f "$CONFIG_DIR/staging.env" ]; then
            if [ -f "$CONFIG_DIR/staging.env.example" ]; then
                cp "$CONFIG_DIR/staging.env.example" "$CONFIG_DIR/staging.env"
                echo "✅ $CONFIG_DIR/staging.env を作成しました"
            else
                echo "⚠️  $CONFIG_DIR/staging.env.example が見つかりません"
            fi
        else
            echo "ℹ️  $CONFIG_DIR/staging.env は既に存在します"
        fi
        ;;
    production)
        if [ ! -f "$CONFIG_DIR/production.env" ]; then
            if [ -f "$CONFIG_DIR/production.env.example" ]; then
                cp "$CONFIG_DIR/production.env.example" "$CONFIG_DIR/production.env"
                echo "✅ $CONFIG_DIR/production.env を作成しました"
            else
                echo "⚠️  $CONFIG_DIR/production.env.example が見つかりません"
            fi
        else
            echo "ℹ️  $CONFIG_DIR/production.env は既に存在します"
        fi
        ;;
    *)
        echo "❌ エラー: 無効な環境です。development, staging, production のいずれかを指定してください"
        exit 1
        ;;
esac

echo ""
echo "📝 次のステップ:"
echo "   1. $CONFIG_DIR/$ENVIRONMENT.env を編集して、実際の設定値を入力してください"
echo "   2. 本番環境の場合は、機密情報を安全に管理してください"
echo ""

