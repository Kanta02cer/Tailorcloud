package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tailor-cloud/backend/internal/middleware"
	"tailor-cloud/backend/internal/service"
)

// CustomerHandler 顧客ハンドラー
type CustomerHandler struct {
	customerService *service.CustomerService
}

// NewCustomerHandler CustomerHandlerのコンストラクタ
func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomerRequest 顧客作成リクエスト
type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// CreateCustomer POST /api/customers - 顧客を登録
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCustomerRequest
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

	// サービス層で作成
	serviceReq := &service.CreateCustomerRequest{
		TenantID: authUser.TenantID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	customer, err := h.customerService.CreateCustomer(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomer GET /api/customers/{id} - 顧客詳細を取得
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから顧客IDを取得
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		// パスパラメータから取得を試みる（Go 1.22+）
		if id := r.PathValue("id"); id != "" {
			customerID = id
		}
	}
	if customerID == "" {
		// フォールバック: パスから手動パース
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "customers" {
			customerID = pathParts[len(pathParts)-1]
		}
	}

	if customerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
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

	// サービス層で取得
	customer, err := h.customerService.GetCustomer(r.Context(), customerID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get customer: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// ListCustomers GET /api/customers - 顧客一覧を取得（検索対応）
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
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

	tenantID := authUser.TenantID
	searchKeyword := r.URL.Query().Get("search")

	// サービス層で一覧取得
	serviceReq := &service.ListCustomersRequest{
		TenantID: tenantID,
		Search:   searchKeyword,
	}

	customers, err := h.customerService.ListCustomers(r.Context(), serviceReq)
	if err != nil {
		http.Error(w, "Failed to list customers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンス
	response := map[string]interface{}{
		"customers": customers,
		"total":     len(customers),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateCustomerRequest 顧客更新リクエスト
type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UpdateCustomer PUT /api/customers/{id} - 顧客を更新
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから顧客IDを取得
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		if id := r.PathValue("id"); id != "" {
			customerID = id
		}
	}
	if customerID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "customers" {
			customerID = pathParts[len(pathParts)-1]
		}
	}

	if customerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	var req UpdateCustomerRequest
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

	// サービス層で更新
	serviceReq := &service.UpdateCustomerRequest{
		TenantID:   authUser.TenantID,
		CustomerID: customerID,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
	}

	customer, err := h.customerService.UpdateCustomer(r.Context(), serviceReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to update customer: "+err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomer DELETE /api/customers/{id} - 顧客を削除
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスまたはクエリパラメータから顧客IDを取得
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		if id := r.PathValue("id"); id != "" {
			customerID = id
		}
	}
	if customerID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 3 && pathParts[len(pathParts)-2] == "customers" {
			customerID = pathParts[len(pathParts)-1]
		}
	}

	if customerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
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

	// サービス層で削除
	if err := h.customerService.DeleteCustomer(r.Context(), customerID, authUser.TenantID); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to delete customer: "+err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCustomerOrders GET /api/customers/{id}/orders - 顧客の注文履歴を取得
func (h *CustomerHandler) GetCustomerOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLパスから顧客IDを取得
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		if id := r.PathValue("id"); id != "" {
			customerID = id
		}
	}
	if customerID == "" {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) >= 4 && pathParts[len(pathParts)-1] == "orders" {
			customerID = pathParts[len(pathParts)-2]
		}
	}

	if customerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
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

	// サービス層で注文履歴を取得
	orders, err := h.customerService.GetCustomerOrders(r.Context(), customerID, authUser.TenantID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		http.Error(w, "Failed to get customer orders: "+err.Error(), statusCode)
		return
	}

	// レスポンス
	response := map[string]interface{}{
		"orders": orders,
		"total":  len(orders),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

