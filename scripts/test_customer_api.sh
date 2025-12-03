#!/bin/bash

# TailorCloud: 顧客管理API動作確認スクリプト

set -e

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== TailorCloud 顧客管理API動作確認 ==="
echo ""

# デフォルト値
API_BASE_URL=${API_BASE_URL:-http://localhost:8080}
TENANT_ID=${TENANT_ID:-tenant-123}

echo -e "${GREEN}✅ 設定:${NC}"
echo "  API_BASE_URL: $API_BASE_URL"
echo "  TENANT_ID: $TENANT_ID"
echo ""

# ヘルスチェック
echo -e "${YELLOW}📋 1. ヘルスチェック${NC}"
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "$API_BASE_URL/health" || echo "000")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
BODY=$(echo "$HEALTH_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ バックエンドAPIは正常に動作しています${NC}"
    echo "  レスポンス: $BODY"
else
    echo -e "${RED}❌ バックエンドAPIに接続できません${NC}"
    echo "  HTTPステータス: $HTTP_CODE"
    echo "  レスポンス: $BODY"
    echo ""
    echo "バックエンドAPIを起動してください:"
    echo "  ./scripts/start_backend.sh"
    exit 1
fi

echo ""

# 顧客一覧取得
echo -e "${YELLOW}📋 2. 顧客一覧取得${NC}"
CUSTOMERS_RESPONSE=$(curl -s -w "\n%{http_code}" "$API_BASE_URL/api/customers?tenant_id=$TENANT_ID" || echo "000")
HTTP_CODE=$(echo "$CUSTOMERS_RESPONSE" | tail -n1)
BODY=$(echo "$CUSTOMERS_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ 顧客一覧取得成功${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
else
    echo -e "${YELLOW}⚠️  顧客一覧取得: HTTP $HTTP_CODE${NC}"
    echo "  レスポンス: $BODY"
    echo "  （認証が必要な場合、エラーが返る可能性があります）"
fi

echo ""

# 顧客作成（テスト用）
echo -e "${YELLOW}📋 3. 顧客作成テスト${NC}"
CUSTOMER_DATA='{
  "name": "テスト顧客",
  "email": "test@example.com",
  "phone": "090-1234-5678",
  "address": "東京都渋谷区テスト1-2-3"
}'

CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$CUSTOMER_DATA" \
  "$API_BASE_URL/api/customers?tenant_id=$TENANT_ID" || echo "000")
HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
BODY=$(echo "$CREATE_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ 顧客作成成功${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
    
    # 作成された顧客IDを抽出（簡易版）
    CUSTOMER_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4 || echo "")
    if [ -n "$CUSTOMER_ID" ]; then
        echo ""
        echo -e "${YELLOW}📋 4. 顧客詳細取得（ID: $CUSTOMER_ID）${NC}"
        GET_RESPONSE=$(curl -s -w "\n%{http_code}" "$API_BASE_URL/api/customers/$CUSTOMER_ID?tenant_id=$TENANT_ID" || echo "000")
        GET_HTTP_CODE=$(echo "$GET_RESPONSE" | tail -n1)
        GET_BODY=$(echo "$GET_RESPONSE" | head -n-1)
        
        if [ "$GET_HTTP_CODE" = "200" ]; then
            echo -e "${GREEN}✅ 顧客詳細取得成功${NC}"
            echo "$GET_BODY" | python3 -m json.tool 2>/dev/null || echo "$GET_BODY"
        else
            echo -e "${YELLOW}⚠️  顧客詳細取得: HTTP $GET_HTTP_CODE${NC}"
            echo "  レスポンス: $GET_BODY"
        fi
    fi
else
    echo -e "${YELLOW}⚠️  顧客作成: HTTP $HTTP_CODE${NC}"
    echo "  レスポンス: $BODY"
    echo "  （認証が必要な場合、エラーが返る可能性があります）"
fi

echo ""
echo -e "${GREEN}✅ 動作確認完了${NC}"
echo ""
echo "注意: 認証が必要なAPIは、Firebase認証トークンが必要です。"
echo "実際の動作確認は、Flutterアプリから行ってください。"

