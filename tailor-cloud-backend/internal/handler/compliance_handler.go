package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// ComplianceHandler コンプライアンスハンドラー
type ComplianceHandler struct {
	complianceService *service.ComplianceService
	orderService      *service.OrderService
}

// NewComplianceHandler ComplianceHandlerのコンストラクタ
func NewComplianceHandler(
	complianceService *service.ComplianceService,
	orderService *service.OrderService,
) *ComplianceHandler {
	return &ComplianceHandler{
		complianceService: complianceService,
		orderService:      orderService,
	}
}

// GenerateDocumentRequest 発注書生成リクエスト
type GenerateDocumentRequest struct {
	// リクエストボディは空でもOK（注文IDから必要な情報を取得）
}

// GenerateDocumentResponse 発注書生成レスポンス
type GenerateDocumentResponse struct {
	OrderID     string `json:"order_id"`
	DocURL      string `json:"doc_url"`
	DocHash     string `json:"doc_hash"`
	GeneratedAt string `json:"generated_at"`
}

// GenerateDocument POST /api/orders/{id}/generate-document - 発注書PDFを生成
func (h *ComplianceHandler) GenerateDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// URLパスから注文IDを取得
	// パターン1: /api/orders/{id}/generate-document (Go 1.22+のパスパラメータ)
	// パターン2: /api/orders/generate-document?order_id={id} (クエリパラメータ、フォールバック)
	
	var orderID string
	
	// Go 1.22+のPathValueメソッドでパスパラメータから取得を試みる
	// PathValueはパスパラメータが存在しない場合は空文字列を返す
	if id := r.PathValue("id"); id != "" {
		orderID = id
	}
	
	// パスパラメータで取得できなかった場合はクエリパラメータから取得
	if orderID == "" {
		orderID = r.URL.Query().Get("order_id")
	}
	
	// それでも取得できない場合はパスから手動パース（フォールバック）
	if orderID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 4 && pathParts[len(pathParts)-1] == "generate-document" {
			orderID = pathParts[len(pathParts)-2]
		}
	}
	
	if orderID == "" {
		http.Error(w, "order_id is required (path parameter or query parameter)", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 注文を取得
	order, err := h.orderService.GetOrder(ctx, orderID, authUser.TenantID)
	if err != nil {
		http.Error(w, "Order not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// テナント情報を構築（TODO: Tenantリポジトリから取得）
	// MVPでは簡易的に注文情報から構築
	tenant := &domain.Tenant{
		ID:        authUser.TenantID,
		LegalName: "Regalis Group", // TODO: Tenantリポジトリから取得
		Type:      domain.TenantTypeTailor,
	}

	// コンプライアンス要件を構築
	requirement := domain.BuildComplianceRequirementFromOrder(order, tenant, order.Details)
	if requirement == nil {
		http.Error(w, "Failed to build compliance requirement", http.StatusInternalServerError)
		return
	}

	// PDF生成リクエストを作成
	pdfReq := &service.GenerateComplianceDocumentRequest{
		Order:       order,
		Tenant:      tenant,
		Requirement: requirement,
	}

	// PDF生成
	pdfResp, err := h.complianceService.GenerateComplianceDocument(ctx, pdfReq)
	if err != nil {
		http.Error(w, "Failed to generate PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 注文にPDF情報を更新
	// TODO: OrderServiceにUpdateComplianceDocメソッドを追加
	// 現時点ではレスポンスのみ返す

	// レスポンスを返す
	response := GenerateDocumentResponse{
		OrderID:     order.ID,
		DocURL:      pdfResp.DocURL,
		DocHash:     pdfResp.DocHash,
		GeneratedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GenerateAmendmentDocumentRequest 修正発注書生成リクエスト
type GenerateAmendmentDocumentRequest struct {
	AmendmentReason string `json:"amendment_reason"` // 修正理由（必須）
}

// GenerateAmendmentDocument POST /api/orders/{id}/generate-amendment - 修正発注書PDFを生成
func (h *ComplianceHandler) GenerateAmendmentDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// URLパスから注文IDを取得
	var orderID string
	if id := r.PathValue("id"); id != "" {
		orderID = id
	}
	if orderID == "" {
		orderID = r.URL.Query().Get("order_id")
	}
	if orderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	// リクエストボディから修正理由を取得
	var reqBody GenerateAmendmentDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.AmendmentReason == "" {
		http.Error(w, "amendment_reason is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 注文を取得
	order, err := h.orderService.GetOrder(ctx, orderID, authUser.TenantID)
	if err != nil {
		http.Error(w, "Order not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// テナント情報を構築（TODO: Tenantリポジトリから取得）
	tenant := &domain.Tenant{
		ID:        authUser.TenantID,
		LegalName: "Regalis Group", // TODO: Tenantリポジトリから取得
		Type:      domain.TenantTypeTailor,
	}

	// 修正発注書生成リクエスト
	// GeneratedByにはユーザーIDを設定（AuthUser.IDを使用）
	generatedBy := authUser.ID
	if generatedBy == "" {
		generatedBy = "system" // デフォルト値（認証情報がない場合）
	}
	
	serviceReq := &service.GenerateAmendmentDocumentRequest{
		OrderID:         orderID,
		TenantID:        authUser.TenantID,
		GeneratedBy:     generatedBy,
		AmendmentReason: reqBody.AmendmentReason,
	}

	// 修正発注書を生成
	resp, err := h.complianceService.GenerateAmendmentDocument(ctx, serviceReq, order, tenant)
	if err != nil {
		http.Error(w, "Failed to generate amendment document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetComplianceDocuments GET /api/orders/{id}/compliance-documents - 発注書履歴を取得
func (h *ComplianceHandler) GetComplianceDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// URLパスから注文IDを取得
	var orderID string
	if id := r.PathValue("id"); id != "" {
		orderID = id
	}
	if orderID == "" {
		orderID = r.URL.Query().Get("order_id")
	}
	if orderID == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザー情報をコンテキストから取得
	authUser, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// 注文を取得（存在確認）
	_, err = h.orderService.GetOrder(ctx, orderID, authUser.TenantID)
	if err != nil {
		http.Error(w, "Order not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// TODO: ComplianceDocumentRepositoryをハンドラーに注入する必要がある
	// 現時点では、ComplianceServiceにGetDocumentsメソッドを追加するか、
	// 直接リポジトリを使用する必要がある
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}
