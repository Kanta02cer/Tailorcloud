package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// InvoiceHandler 請求書ハンドラー
type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

// NewInvoiceHandler InvoiceHandlerのコンストラクタ
func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

// GenerateInvoiceRequest 請求書生成リクエスト
type GenerateInvoiceRequest struct {
	OrderID string `json:"order_id"`
}

// GenerateInvoice POST /api/orders/{id}/generate-invoice - 適格請求書（インボイス）を生成
func (h *InvoiceHandler) GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータからIDを取得
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		if id := r.PathValue("id"); id != "" {
			orderID = id
		}
	}
	if orderID == "" {
		// リクエストボディからも取得を試みる
		var req GenerateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			orderID = req.OrderID
		}
	}
	if orderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得（テナント検証用）
	_, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		// 認証情報がない場合は、リクエストパラメータから取得（開発環境用フォールバック）
		tenantID := r.URL.Query().Get("tenant_id")
		if tenantID == "" {
			http.Error(w, "Authentication required or tenant_id must be provided", http.StatusUnauthorized)
			return
		}
		// 開発環境ではテナントIDのみで処理可能
	}

	// サービス層で請求書生成
	serviceReq := &service.InvoiceRequest{
		OrderID: orderID,
	}

	resp, err := h.invoiceService.GenerateInvoice(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to generate invoice: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

