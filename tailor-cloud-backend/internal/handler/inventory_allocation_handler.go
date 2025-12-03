package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// InventoryAllocationHandler 在庫引当ハンドラー
type InventoryAllocationHandler struct {
	allocationService *service.InventoryAllocationService
}

// NewInventoryAllocationHandler InventoryAllocationHandlerのコンストラクタ
func NewInventoryAllocationHandler(allocationService *service.InventoryAllocationService) *InventoryAllocationHandler {
	return &InventoryAllocationHandler{
		allocationService: allocationService,
	}
}

// AllocateInventoryRequest 在庫引当リクエスト
type AllocateInventoryRequest struct {
	OrderID        string  `json:"order_id"`
	FabricID       string  `json:"fabric_id"`
	RequiredLength float64 `json:"required_length"` // 必要な長さ（メートル）
	Strategy       string  `json:"strategy,omitempty"` // FIFO, LIFO, BEST_FIT
}

// AllocateInventory POST /api/inventory/allocate - 在庫を引当
func (h *InventoryAllocationHandler) AllocateInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AllocateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}
	if req.FabricID == "" {
		http.Error(w, "fabric_id is required", http.StatusBadRequest)
		return
	}
	if req.RequiredLength <= 0 {
		http.Error(w, "required_length must be greater than 0", http.StatusBadRequest)
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

	// 引当戦略のパース（デフォルトはFIFO）
	strategy := service.AllocationStrategyFIFO
	if req.Strategy != "" {
		strategy = service.AllocationStrategy(req.Strategy)
	}

	// サービス層で引当
	serviceReq := &service.AllocateInventoryRequest{
		TenantID:       authUser.TenantID,
		OrderID:        req.OrderID,
		FabricID:       req.FabricID,
		RequiredLength: req.RequiredLength,
		Strategy:       strategy,
	}

	resp, err := h.allocationService.AllocateInventory(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "insufficient inventory") {
			statusCode = http.StatusBadRequest
		}
		http.Error(w, "Failed to allocate inventory: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// ReleaseAllocationRequest 引当解除リクエスト
type ReleaseAllocationRequest struct {
	AllocationID string `json:"allocation_id"`
}

// ReleaseAllocation POST /api/inventory/release - 引当を解除（キャンセル時など）
func (h *InventoryAllocationHandler) ReleaseAllocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReleaseAllocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AllocationID == "" {
		http.Error(w, "allocation_id is required", http.StatusBadRequest)
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

	// サービス層で解除
	if err := h.allocationService.ReleaseAllocation(r.Context(), req.AllocationID, authUser.TenantID); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to release allocation: "+err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

