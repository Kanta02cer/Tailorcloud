#!/bin/bash

# TailorCloud バックエンド デプロイスクリプト
# Google Cloud Runへのデプロイ

set -e

echo "🚀 TailorCloud バックエンドをデプロイします..."

# プロジェクト設定
PROJECT_ID=${GCP_PROJECT_ID:-"gen-lang-client-0552849356"}
REGION=${REGION:-"asia-northeast1"}
SERVICE_NAME="tailor-cloud-backend"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "📋 設定:"
echo "  プロジェクトID: ${PROJECT_ID}"
echo "  リージョン: ${REGION}"
echo "  サービス名: ${SERVICE_NAME}"
echo ""

# バックエンドディレクトリに移動
cd "$(dirname "$0")/../tailor-cloud-backend"

# Dockerイメージをビルド
echo "🔨 Dockerイメージをビルド中..."
gcloud builds submit --tag ${IMAGE_NAME} --project ${PROJECT_ID}

# Cloud Runにデプロイ
echo "☁️  Cloud Runにデプロイ中..."
gcloud run deploy ${SERVICE_NAME} \
  --image ${IMAGE_NAME} \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --set-env-vars "PORT=8080" \
  --set-env-vars "GCP_PROJECT_ID=${PROJECT_ID}" \
  --project ${PROJECT_ID}

# デプロイされたURLを取得
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} \
  --platform managed \
  --region ${REGION} \
  --format 'value(status.url)' \
  --project ${PROJECT_ID})

echo ""
echo "✅ デプロイ完了!"
echo "🌐 サービスURL: ${SERVICE_URL}"
echo ""
echo "📊 ヘルスチェック:"
echo "  curl ${SERVICE_URL}/health"
echo ""
