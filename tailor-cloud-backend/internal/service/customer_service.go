package service

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// CustomerService 顧客サービス
type CustomerService struct {
	customerRepo repository.CustomerRepository
	orderRepo    repository.OrderRepository
}

// NewCustomerService CustomerServiceのコンストラクタ
func NewCustomerService(customerRepo repository.CustomerRepository, orderRepo repository.OrderRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		orderRepo:    orderRepo,
	}
}

// CreateCustomerRequest 顧客作成リクエスト
type CreateCustomerRequest struct {
	TenantID string
	Name     string
	Email    string
	Phone    string
}

// CreateCustomer 顧客を作成
func (s *CustomerService) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*domain.Customer, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	
	customer := &domain.Customer{
		TenantID: req.TenantID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	
	return customer, nil
}

// GetCustomer 顧客を取得
func (s *CustomerService) GetCustomer(ctx context.Context, customerID string, tenantID string) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, customerID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	return customer, nil
}

// ListCustomersRequest 顧客一覧取得リクエスト
type ListCustomersRequest struct {
	TenantID string
	Search   string // 検索キーワード（名前、メール、電話番号）
}

// ListCustomers 顧客一覧を取得
func (s *CustomerService) ListCustomers(ctx context.Context, req *ListCustomersRequest) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	var err error
	
	if req.Search != "" {
		// 検索あり
		customers, err = s.customerRepo.Search(ctx, req.TenantID, req.Search)
	} else {
		// 検索なし（全件取得）
		customers, err = s.customerRepo.GetByTenantID(ctx, req.TenantID)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	
	return customers, nil
}

// UpdateCustomerRequest 顧客更新リクエスト
type UpdateCustomerRequest struct {
	TenantID string
	CustomerID string
	Name     string
	Email    string
	Phone    string
}

// UpdateCustomer 顧客を更新
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *UpdateCustomerRequest) (*domain.Customer, error) {
	if req.CustomerID == "" {
		return nil, fmt.Errorf("customer_id is required")
	}
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	
	// 既存の顧客を取得
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	// 更新
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}
	
	return customer, nil
}

// DeleteCustomer 顧客を削除
func (s *CustomerService) DeleteCustomer(ctx context.Context, customerID string, tenantID string) error {
	if customerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	
	if err := s.customerRepo.Delete(ctx, customerID, tenantID); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	
	return nil
}

// GetCustomerOrders 顧客の注文履歴を取得
func (s *CustomerService) GetCustomerOrders(ctx context.Context, customerID string, tenantID string) ([]*domain.Order, error) {
	// 顧客が存在するか確認
	_, err := s.customerRepo.GetByID(ctx, customerID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	
	// テナントの全注文を取得
	allOrders, err := s.orderRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	
	// 顧客IDでフィルター
	var customerOrders []*domain.Order
	for _, order := range allOrders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	
	return customerOrders, nil
}

