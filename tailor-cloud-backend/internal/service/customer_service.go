package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

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
type CustomerInteractionInput struct {
	Type      string     `json:"type"`
	Note      string     `json:"note"`
	Staff     string     `json:"staff"`
	Timestamp *time.Time `json:"timestamp"`
}

// CreateCustomerRequest 顧客作成リクエスト
type CreateCustomerRequest struct {
	TenantID           string
	Name               string
	Email              string
	Phone              string
	Address            string
	CustomerStatus     string
	Tags               []string
	VIPRank            int
	LTVScore           float64
	LifetimeValue      float64
	PreferredChannel   string
	LeadSource         string
	Notes              string
	Occupation         string
	AnnualIncomeRange  string
	PreferredArchetype string
	DiagnosisCount     int
	LastInteractionAt  *time.Time
	Interactions       []CustomerInteractionInput
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
		TenantID:           req.TenantID,
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Phone,
		Address:            strings.TrimSpace(req.Address),
		CustomerStatus:     normalizeStatus(req.CustomerStatus),
		Tags:               normalizeTags(req.Tags),
		VIPRank:            req.VIPRank,
		LTVScore:           req.LTVScore,
		LifetimeValue:      req.LifetimeValue,
		PreferredChannel:   strings.TrimSpace(req.PreferredChannel),
		LeadSource:         strings.TrimSpace(req.LeadSource),
		Notes:              req.Notes,
		Occupation:         req.Occupation,
		AnnualIncomeRange:  req.AnnualIncomeRange,
		PreferredArchetype: req.PreferredArchetype,
		DiagnosisCount:     req.DiagnosisCount,
	}

	customer.InteractionNotes = convertInteractionInputs(req.Interactions)
	customer.LastInteractionAt = determineLastInteraction(req.LastInteractionAt, customer.InteractionNotes)

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
	TenantID           string
	CustomerID         string
	Name               string
	Email              string
	Phone              string
	Address            string
	CustomerStatus     string
	Tags               []string
	VIPRank            int
	LTVScore           float64
	LifetimeValue      float64
	PreferredChannel   string
	LeadSource         string
	Notes              string
	Occupation         string
	AnnualIncomeRange  string
	PreferredArchetype string
	DiagnosisCount     int
	LastInteractionAt  *time.Time
	Interactions       []CustomerInteractionInput
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
	customer.Address = strings.TrimSpace(req.Address)
	customer.CustomerStatus = normalizeStatus(req.CustomerStatus)
	customer.Tags = normalizeTags(req.Tags)
	customer.VIPRank = req.VIPRank
	customer.LTVScore = req.LTVScore
	customer.LifetimeValue = req.LifetimeValue
	customer.PreferredChannel = strings.TrimSpace(req.PreferredChannel)
	customer.LeadSource = strings.TrimSpace(req.LeadSource)
	customer.Notes = req.Notes
	customer.Occupation = req.Occupation
	customer.AnnualIncomeRange = req.AnnualIncomeRange
	customer.PreferredArchetype = req.PreferredArchetype
	customer.DiagnosisCount = req.DiagnosisCount

	if len(req.Interactions) > 0 {
		customer.InteractionNotes = convertInteractionInputs(req.Interactions)
	}
	customer.LastInteractionAt = determineLastInteraction(req.LastInteractionAt, customer.InteractionNotes)

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

func normalizeStatus(status string) string {
	if status == "" {
		return "lead"
	}
	switch strings.ToLower(status) {
	case "lead", "prospect", "active", "inactive", "vip":
		return strings.ToLower(status)
	default:
		return "lead"
	}
}

func normalizeTags(tags []string) []string {
	dedup := make(map[string]struct{})
	for _, tag := range tags {
		trimmed := strings.TrimSpace(strings.ToLower(tag))
		if trimmed == "" {
			continue
		}
		dedup[trimmed] = struct{}{}
	}
	if len(dedup) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(dedup))
	for tag := range dedup {
		result = append(result, tag)
	}
	sort.Strings(result)
	return result
}

func convertInteractionInputs(inputs []CustomerInteractionInput) []domain.CustomerInteractionNote {
	if len(inputs) == 0 {
		return []domain.CustomerInteractionNote{}
	}
	notes := make([]domain.CustomerInteractionNote, 0, len(inputs))
	for _, input := range inputs {
		ts := time.Now()
		if input.Timestamp != nil && !input.Timestamp.IsZero() {
			ts = *input.Timestamp
		}
		noteType := input.Type
		if noteType == "" {
			noteType = "note"
		}
		notes = append(notes, domain.CustomerInteractionNote{
			Timestamp: ts,
			Type:      strings.ToLower(noteType),
			Note:      input.Note,
			Staff:     input.Staff,
		})
	}
	return notes
}

func determineLastInteraction(explicit *time.Time, notes []domain.CustomerInteractionNote) *time.Time {
	if explicit != nil && !explicit.IsZero() {
		return explicit
	}
	var latest *time.Time
	for _, note := range notes {
		if latest == nil || note.Timestamp.After(*latest) {
			t := note.Timestamp
			latest = &t
		}
	}
	return latest
}
