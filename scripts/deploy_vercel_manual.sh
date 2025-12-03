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

# プロジェクトのリンク確認
if [ ! -f .vercel/project.json ]; then
    echo "🔗 プロジェクトをリンク中..."
    vercel link --project tailorcloud --yes
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
echo "  3) Vercelダッシュボードからデプロイ（推奨 - Git authorエラー回避）"
read -p "選択 (1, 2, or 3): " deploy_type

case $deploy_type in
    1)
        echo ""
        echo "🚀 プレビューデプロイを開始..."
        vercel
        ;;
    2)
        echo ""
        echo "🚀 本番環境デプロイを開始..."
        echo "⚠️  Git authorエラーが発生する場合は、オプション3を選択してください"
        vercel --prod || {
            echo ""
            echo "❌ デプロイエラーが発生しました"
            echo "💡 解決方法: Vercelダッシュボードから再デプロイしてください"
            echo "   URL: https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud"
            exit 1
        }
        ;;
    3)
        echo ""
        echo "🌐 Vercelダッシュボードを開きます..."
        echo "   1. プロジェクトを選択"
        echo "   2. 'Deployments' タブをクリック"
        echo "   3. 最新のデプロイメントの '...' メニュー → 'Redeploy' をクリック"
        echo ""
        open "https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud" 2>/dev/null || {
            echo "   ブラウザで以下のURLにアクセスしてください:"
            echo "   https://vercel.com/kinouecertify-gmailcoms-projects/tailorcloud"
        }
        ;;
    *)
        echo "❌ 無効な選択です"
        exit 1
        ;;
esac

echo ""
echo "=== デプロイ完了 ==="

