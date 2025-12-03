package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// AuthContextKey コンテキストキー（型安全なキー）
type AuthContextKey string

const (
	UserIDKey   AuthContextKey = "user_id"
	TenantIDKey AuthContextKey = "tenant_id"
	RoleKey     AuthContextKey = "role"
)

// AuthUser 認証済みユーザー情報
type AuthUser struct {
	ID       string
	TenantID string
	Role     string
	Email    string
}

// FirebaseAuthMiddleware Firebase認証ミドルウェア
// JWTトークンを検証し、ユーザー情報をコンテキストに注入
type FirebaseAuthMiddleware struct {
	authClient *auth.Client
}

// NewFirebaseAuthMiddleware FirebaseAuthMiddlewareのコンストラクタ
func NewFirebaseAuthMiddleware(app *firebase.App) (*FirebaseAuthMiddleware, error) {
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth client: %w", err)
	}

	return &FirebaseAuthMiddleware{
		authClient: authClient,
	}, nil
}

// Authenticate 認証ミドルウェアハンドラー
func (m *FirebaseAuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーからトークンを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// "Bearer <token>" 形式をパース
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format. Expected: Bearer <token>", http.StatusUnauthorized)
			return
		}

		idToken := parts[1]

		// Firebase Authでトークンを検証
		ctx := r.Context()
		token, err := m.authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// カスタムクレームからユーザー情報を取得
		// Firebase Authのカスタムクレームで tenant_id と role を設定する前提
		userID := token.UID
		
		// カスタムクレームからテナントIDとロールを取得
		tenantID, _ := token.Claims["tenant_id"].(string)
		role, _ := token.Claims["role"].(string)

		// ユーザー情報をコンテキストに注入
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, TenantIDKey, tenantID)
		ctx = context.WithValue(ctx, RoleKey, role)

		// リクエストのコンテキストを更新
		r = r.WithContext(ctx)

		// 次のハンドラーへ
		next.ServeHTTP(w, r)
	}
}

// GetUserFromContext コンテキストからユーザー情報を取得
func GetUserFromContext(ctx context.Context) (*AuthUser, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id not found in context")
	}

	tenantID, _ := ctx.Value(TenantIDKey).(string)
	role, _ := ctx.Value(RoleKey).(string)

	return &AuthUser{
		ID:       userID,
		TenantID: tenantID,
		Role:     role,
	}, nil
}

// OptionalAuth オプショナル認証ミドルウェア
// 認証が失敗してもリクエストを通す（開発環境用）
func (m *FirebaseAuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// トークンがない場合はそのまま通す
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// トークンがある場合は検証を試みる
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 形式が不正でも通す（開発環境用）
			next.ServeHTTP(w, r)
			return
		}

		idToken := parts[1]
		ctx := r.Context()

		// 検証を試みる（失敗してもエラーを返さない）
		token, err := m.authClient.VerifyIDToken(ctx, idToken)
		if err == nil {
			// 検証成功時のみコンテキストに注入
			userID := token.UID
			tenantID, _ := token.Claims["tenant_id"].(string)
			role, _ := token.Claims["role"].(string)

			ctx = context.WithValue(ctx, UserIDKey, userID)
			ctx = context.WithValue(ctx, TenantIDKey, tenantID)
			ctx = context.WithValue(ctx, RoleKey, role)

			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRole ロールチェックミドルウェア
func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// 許可されたロールかチェック
			for _, allowedRole := range allowedRoles {
				if user.Role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			// 権限不足
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Insufficient permissions. Required role: %v", allowedRoles),
			})
		}
	}
}

