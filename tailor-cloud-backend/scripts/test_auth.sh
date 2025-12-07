#!/bin/bash

# 認証エンドポイントのテストスクリプト
# 使用方法: ./scripts/test_auth.sh <firebase-id-token>

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
ID_TOKEN="${1:-}"

if [ -z "$ID_TOKEN" ]; then
    echo "Error: Firebase ID token is required"
    echo "Usage: ./scripts/test_auth.sh <firebase-id-token>"
    echo ""
    echo "Example:"
    echo "  ./scripts/test_auth.sh 'eyJhbGciOiJSUzI1NiIsImtpZCI6...'"
    exit 1
fi

echo "Testing authentication endpoint..."
echo "API Base URL: $API_BASE_URL"
echo ""

# 認証エンドポイントをテスト
response=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE_URL/api/auth/verify" \
  -H "Content-Type: application/json" \
  -d "{\"id_token\": \"$ID_TOKEN\"}")

# レスポンスとステータスコードを分離
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "HTTP Status Code: $http_code"
echo "Response Body:"
echo "$body" | jq '.' 2>/dev/null || echo "$body"
echo ""

if [ "$http_code" -eq 200 ]; then
    echo "✅ Authentication successful!"
    verified=$(echo "$body" | jq -r '.verified' 2>/dev/null || echo "false")
    if [ "$verified" = "true" ]; then
        user_id=$(echo "$body" | jq -r '.user_id' 2>/dev/null || echo "N/A")
        email=$(echo "$body" | jq -r '.email' 2>/dev/null || echo "N/A")
        tenant_id=$(echo "$body" | jq -r '.tenant_id' 2>/dev/null || echo "N/A")
        role=$(echo "$body" | jq -r '.role' 2>/dev/null || echo "N/A")
        
        echo ""
        echo "User Information:"
        echo "  User ID: $user_id"
        echo "  Email: $email"
        echo "  Tenant ID: $tenant_id"
        echo "  Role: $role"
    fi
else
    echo "❌ Authentication failed!"
    exit 1
fi

