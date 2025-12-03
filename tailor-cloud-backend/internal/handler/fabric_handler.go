package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// FabricHandler 生地ハンドラー
type FabricHandler struct {
	fabricService *service.FabricService
}

// NewFabricHandler FabricHandlerのコンストラクタ
func NewFabricHandler(fabricService *service.FabricService) *FabricHandler {
	return &FabricHandler{
		fabricService: fabricService,
	}
}

// ListFabrics GET /api/fabrics - 生地一覧を取得（フィルター・検索対応）
func (h *FabricHandler) ListFabrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		// 認証情報がない場合は、リクエストパラメータから取得（開発環境用フォールバック）
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}
	
	// テナントID: 認証ユーザーから取得、またはクエリパラメータから
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}
	
	// フィルター・検索パラメータ
	statusParam := r.URL.Query().Get("status")
	searchParam := r.URL.Query().Get("search")
	
	// ステータスフィルターをパース
	var statusFilter []domain.StockStatus
	if statusParam != "" && statusParam != "all" {
		statuses := strings.Split(statusParam, ",")
		for _, s := range statuses {
			switch strings.TrimSpace(s) {
			case "available":
				statusFilter = append(statusFilter, domain.StockStatusAvailable)
			case "limited":
				statusFilter = append(statusFilter, domain.StockStatusLimited)
			case "soldout":
				statusFilter = append(statusFilter, domain.StockStatusSoldOut)
			}
		}
	}
	
	// サービス層で一覧取得
	req := &service.ListFabricsRequest{
		TenantID: tenantID,
		Status:   statusFilter,
		Search:   searchParam,
	}
	
	fabrics, err := h.fabricService.ListFabrics(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to list fabrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// レスポンス
	response := map[string]interface{}{
		"fabrics": fabrics,
		"total":   len(fabrics),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetFabric GET /api/fabrics/{fabric_id} - 生地詳細を取得
func (h *FabricHandler) GetFabric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// URLからfabric_idを取得
	fabricID := r.URL.Query().Get("fabric_id")
	if fabricID == "" {
		// TODO: パスパラメータから取得する実装（mux/v2等を使用）
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
		authUser = &middleware.AuthUser{TenantID: tenantID}
	}
	
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}
	
	// サービス層で取得
	req := &service.GetFabricRequest{
		FabricID: fabricID,
		TenantID: tenantID,
	}
	
	fabric, err := h.fabricService.GetFabric(r.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get fabric: "+err.Error(), statusCode)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fabric)
}

// ReserveFabricRequest 生地確保リクエスト
type ReserveFabricRequest struct {
	FabricID string  `json:"fabric_id"`
	TenantID string  `json:"tenant_id"`
	Amount   float64 `json:"amount"` // 確保したい数量（メートル）
}

// ReserveFabric POST /api/fabrics/{fabric_id}/reserve - 生地を確保（発注フロー開始）
func (h *FabricHandler) ReserveFabric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req ReserveFabricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		if req.TenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
		authUser = &middleware.AuthUser{TenantID: req.TenantID}
	}
	
	// テナントID: 認証ユーザーから取得、またはリクエストから
	tenantID := authUser.TenantID
	if tenantID == "" {
		tenantID = req.TenantID
		if tenantID == "" {
			http.Error(w, "tenant_id is required", http.StatusBadRequest)
			return
		}
	}
	
	// バリデーション
	if req.FabricID == "" {
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}
	
	// サービス層で確保
	reserveReq := &service.ReserveFabricRequest{
		FabricID: req.FabricID,
		TenantID: tenantID,
		Amount:   req.Amount,
	}
	
	if err := h.fabricService.ReserveFabric(r.Context(), reserveReq); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "sold out") || strings.Contains(err.Error(), "insufficient") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to reserve fabric: "+err.Error(), statusCode)
		return
	}
	
	// レスポンス
	response := map[string]interface{}{
		"message":    "Fabric reservation successful",
		"fabric_id":  req.FabricID,
		"amount":     req.Amount,
		"status":     "reserved",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

