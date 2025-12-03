#!/bin/bash
# TailorCloud: Vercel手動デプロイスクリプト

set -e

echo "=== TailorCloud Vercel手動デプロイ ==="
echo ""

# suit-mbti-web-appディレクトリに移動
cd "$(dirname "$0")/../suit-mbti-web-app" || exit 1

# Vercel CLIがインストールされているか確認
if ! command -v vercel &> /dev/null; then
    echo "❌ Vercel CLIがインストールされていません"
    echo "📦 インストール中..."
    npm install -g vercel
    echo "✅ Vercel CLIインストール完了"
    echo ""
fi

# ログイン状態を確認
if ! vercel whoami &> /dev/null; then
    echo "🔐 Vercelにログインしてください"
    vercel login
    echo ""
fi

# 環境変数の確認
echo "📋 環境変数の確認:"
echo "  VITE_API_BASE_URL: ${VITE_API_BASE_URL:-未設定（デフォルト: http://localhost:8080）}"
echo "  VITE_TENANT_ID: ${VITE_TENANT_ID:-未設定（デフォルト: tenant_test_suit_mbti）}"
echo ""

# デプロイタイプの選択
echo "デプロイタイプを選択してください:"
echo "  1) プレビューデプロイ（開発環境）"
echo "  2) 本番環境デプロイ"
read -p "選択 (1 or 2): " deploy_type

case $deploy_type in
    1)
        echo ""
        echo "🚀 プレビューデプロイを開始..."
        vercel
        ;;
    2)
        echo ""
        echo "🚀 本番環境デプロイを開始..."
        vercel --prod
        ;;
    *)
        echo "❌ 無効な選択です"
        exit 1
        ;;
esac

echo ""
echo "=== デプロイ完了 ==="

