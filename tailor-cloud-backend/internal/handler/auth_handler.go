package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/service"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// AuthHandler 認証ハンドラー
type AuthHandler struct {
	authClient  *auth.Client
	userService *service.UserService
}

// NewAuthHandler AuthHandlerのコンストラクタ
func NewAuthHandler(app *firebase.App, userService *service.UserService) (*AuthHandler, error) {
	ctx := context.Background()
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		authClient:  authClient,
		userService: userService,
	}, nil
}

// VerifyTokenRequest トークン検証リクエスト
type VerifyTokenRequest struct {
	IDToken string `json:"id_token"`
}

// VerifyTokenResponse トークン検証レスポンス
type VerifyTokenResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	TenantID string `json:"tenant_id,omitempty"`
	Role     string `json:"role,omitempty"`
	Verified bool   `json:"verified"`
}

// VerifyToken POST /api/auth/verify - Firebase IDトークンを検証
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.IDToken == "" {
		http.Error(w, "id_token is required", http.StatusBadRequest)
		return
	}

	// Firebase Authでトークンを検証
	ctx := r.Context()
	token, err := h.authClient.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyTokenResponse{
			Verified: false,
		})
		return
	}

	// トークンからユーザー情報を取得
	firebaseUID := token.UID
	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)
	if name == "" {
		// nameクレームがない場合はemailから生成
		if email != "" {
			emailParts := strings.Split(email, "@")
			if len(emailParts) > 0 {
				name = emailParts[0]
			}
		}
	}

	// ユーザーサービスでユーザーを取得または作成
	var user *domain.User
	if h.userService != nil {
		createdUser, err := h.userService.GetOrCreateUser(ctx, &service.GetOrCreateUserRequest{
			FirebaseUID: firebaseUID,
			Email:       email,
			Name:        name,
		})
		if err != nil {
			// ユーザー作成に失敗した場合でも、トークン検証は成功しているので情報を返す
			log.Printf("WARNING: Failed to get or create user: %v", err)
		} else {
			user = createdUser
		}
	}

	// レスポンスを返す
	response := VerifyTokenResponse{
		UserID:   firebaseUID,
		Email:    email,
		Verified: true,
	}

	if user != nil {
		response.UserID = user.ID
		response.TenantID = user.TenantID
		response.Role = string(user.Role)
	} else {
		// ユーザーサービスがない場合や作成に失敗した場合は、トークンから取得
		response.TenantID, _ = token.Claims["tenant_id"].(string)
		response.Role, _ = token.Claims["role"].(string)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
