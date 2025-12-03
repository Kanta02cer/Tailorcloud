package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/repository"
)

// FabricRollHandler 反物（Roll）管理ハンドラー
type FabricRollHandler struct {
	rollRepo repository.FabricRollRepository
}

// NewFabricRollHandler FabricRollHandlerのコンストラクタ
func NewFabricRollHandler(rollRepo repository.FabricRollRepository) *FabricRollHandler {
	return &FabricRollHandler{
		rollRepo: rollRepo,
	}
}

// CreateFabricRollRequest 反物作成リクエスト
type CreateFabricRollRequest struct {
	FabricID      string   `json:"fabric_id"`
	RollNumber    string   `json:"roll_number"`
	InitialLength float64  `json:"initial_length"`
	Width         *float64 `json:"width,omitempty"`
	SupplierLotNo *string  `json:"supplier_lot_no,omitempty"`
	ReceivedAt    *string  `json:"received_at,omitempty"` // ISO 8601形式
	Location      *string  `json:"location,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
}

// CreateFabricRoll POST /api/fabric-rolls - 反物（Roll）を作成
func (h *FabricRollHandler) CreateFabricRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateFabricRollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.FabricID == "" {
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}
	if req.RollNumber == "" {
		http.Error(w, "roll_number is required", http.StatusBadRequest)
		return
	}
	if req.InitialLength <= 0 {
		http.Error(w, "initial_length must be greater than 0", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// リクエストをドメインモデルに変換
	roll := &domain.FabricRoll{
		ID:            uuid.New().String(),
		TenantID:      authUser.TenantID,
		FabricID:      req.FabricID,
		RollNumber:    req.RollNumber,
		InitialLength: req.InitialLength,
		CurrentLength: req.InitialLength, // 初期は初期長さと同じ
		Width:         req.Width,
		SupplierLotNo: req.SupplierLotNo,
		Location:      req.Location,
		Status:        domain.FabricRollStatusAvailable,
		Notes:         req.Notes,
	}

	// 入荷日のパース
	if req.ReceivedAt != nil && *req.ReceivedAt != "" {
		receivedAt, err := time.Parse(time.RFC3339, *req.ReceivedAt)
		if err != nil {
			http.Error(w, "Invalid received_at format (expected ISO 8601): "+err.Error(), http.StatusBadRequest)
			return
		}
		roll.ReceivedAt = &receivedAt
	}

	// 作成
	if err := h.rollRepo.Create(r.Context(), roll); err != nil {
		http.Error(w, "Failed to create fabric roll: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(roll)
}

// GetFabricRoll GET /api/fabric-rolls/{id} - 反物詳細を取得
func (h *FabricRollHandler) GetFabricRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータからIDを取得
	rollID := r.URL.Query().Get("roll_id")
	if rollID == "" {
		if id := r.PathValue("id"); id != "" {
			rollID = id
		}
	}
	if rollID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "fabric-rolls" {
			rollID = pathParts[len(pathParts)-1]
		}
	}

	if rollID == "" {
		http.Error(w, "roll_id is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// 取得
	roll, err := h.rollRepo.GetByID(r.Context(), rollID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get fabric roll: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roll)
}

// ListFabricRolls GET /api/fabric-rolls - 反物一覧を取得
func (h *FabricRollHandler) ListFabricRolls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	fabricID := r.URL.Query().Get("fabric_id")
	if fabricID == "" {
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}

	// ステータスフィルター（オプション）
	var status *domain.FabricRollStatus
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		s := domain.FabricRollStatus(statusParam)
		status = &s
	}

	// 一覧取得
	rolls, err := h.rollRepo.ListByFabricID(r.Context(), authUser.TenantID, fabricID, status)
	if err != nil {
		http.Error(w, "Failed to list fabric rolls: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"rolls": rolls,
		"total": len(rolls),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateFabricRollRequest 反物更新リクエスト
type UpdateFabricRollRequest struct {
	RollNumber    *string  `json:"roll_number,omitempty"`
	Width         *float64 `json:"width,omitempty"`
	Location      *string  `json:"location,omitempty"`
	Status        *string  `json:"status,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
}

// UpdateFabricRoll PUT /api/fabric-rolls/{id} - 反物を更新
func (h *FabricRollHandler) UpdateFabricRoll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータからIDを取得
	rollID := r.URL.Query().Get("roll_id")
	if rollID == "" {
		if id := r.PathValue("id"); id != "" {
			rollID = id
		}
	}
	if rollID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "fabric-rolls" {
			rollID = pathParts[len(pathParts)-1]
		}
	}

	if rollID == "" {
		http.Error(w, "roll_id is required", http.StatusBadRequest)
		return
	}

	var req UpdateFabricRollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}

	// 既存の反物を取得
	roll, err := h.rollRepo.GetByID(r.Context(), rollID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get fabric roll: "+err.Error(), statusCode)
		return
	}

	// 更新
	if req.RollNumber != nil {
		roll.RollNumber = *req.RollNumber
	}
	if req.Width != nil {
		roll.Width = req.Width
	}
	if req.Location != nil {
		roll.Location = req.Location
	}
	if req.Status != nil {
		roll.Status = domain.FabricRollStatus(*req.Status)
	}
	if req.Notes != nil {
		roll.Notes = req.Notes
	}

	if err := h.rollRepo.Update(r.Context(), roll); err != nil {
		http.Error(w, "Failed to update fabric roll: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roll)
}

