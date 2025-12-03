package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// OrderService 注文サービス
// ビジネスロジックを実装
type OrderService struct {
	orderRepo         repository.OrderRepository
	auditLogRepo      repository.AuditLogRepository // 監査ログリポジトリ（オプショナル）
	ambassadorService *AmbassadorService            // アンバサダーサービス（成果報酬管理用）
}

// NewOrderService OrderServiceのコンストラクタ
func NewOrderService(orderRepo repository.OrderRepository, auditLogRepo repository.AuditLogRepository, ambassadorService *AmbassadorService) *OrderService {
	return &OrderService{
		orderRepo:         orderRepo,
		auditLogRepo:      auditLogRepo,
		ambassadorService: ambassadorService,
	}
}

// CreateOrderRequest 注文作成リクエスト
type CreateOrderRequest struct {
	TenantID     string               `json:"tenant_id"`
	CustomerID   string               `json:"customer_id"`
	FabricID     string               `json:"fabric_id"`
	TotalAmount  int64                `json:"total_amount"`
	DeliveryDate time.Time            `json:"delivery_date"`
	Details      *domain.OrderDetails `json:"details"`
	CreatedBy    string               `json:"created_by"`
	IPAddress    string               `json:"-"` // HTTPリクエストから取得
	UserAgent    string               `json:"-"` // HTTPリクエストから取得
}

// CreateOrder 注文を作成（Draftステータス）
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*domain.Order, error) {
	// 1. バリデーション
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if req.CustomerID == "" {
		return nil, fmt.Errorf("customer_id is required")
	}
	if req.FabricID == "" {
		return nil, fmt.Errorf("fabric_id is required")
	}
	if req.TotalAmount <= 0 {
		return nil, fmt.Errorf("total_amount must be greater than 0")
	}
	if req.DeliveryDate.IsZero() {
		return nil, fmt.Errorf("delivery_date is required")
	}
	if req.CreatedBy == "" {
		return nil, fmt.Errorf("created_by is required")
	}

	// 2. 注文オブジェクトを作成（Draftステータス）
	order := domain.NewOrder(
		req.TenantID,
		req.CustomerID,
		req.FabricID,
		req.CreatedBy,
		req.TotalAmount,
		req.DeliveryDate,
	)

	// 詳細情報を設定
	if req.Details != nil {
		order.Details = req.Details
	} else {
		order.Details = &domain.OrderDetails{}
	}

	// 3. リポジトリに保存
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 4. 監査ログ記録（非同期・エラー時も継続）
	if s.auditLogRepo != nil {
		var ctxData *auditLogContext = &auditLogContext{
			TenantID:      req.TenantID,
			UserID:        req.CreatedBy,
			Action:        domain.AuditActionCreate,
			ResourceType:  "order",
			ResourceID:    order.ID,
			OldValue:      "",
			NewValue:      s.orderToJSON(order),
			ChangedFields: []string{"all"},
			IPAddress:     req.IPAddress,
			UserAgent:     req.UserAgent,
		}
		s.recordAuditLog(ctxData)
	}

	// 5. アンバサダー成果報酬の作成（非同期・エラー時も継続）
	if s.ambassadorService != nil {
		go func() {
			_, err := s.ambassadorService.CreateCommissionForOrder(context.Background(), order)
			if err != nil {
				fmt.Printf("WARNING: Failed to create commission for order: %v\n", err)
			}
		}()
	}

	return order, nil
}

// ConfirmOrderRequest 注文確定リクエスト
type ConfirmOrderRequest struct {
	OrderID       string `json:"order_id"`
	TenantID      string `json:"tenant_id"`      // セキュリティ: テナントIDを確認
	PrincipalName string `json:"principal_name"` // 委託をする者の氏名
	UserID        string `json:"-"`              // HTTPリクエストから取得
	IPAddress     string `json:"-"`              // HTTPリクエストから取得
	UserAgent     string `json:"-"`              // HTTPリクエストから取得
}

// ConfirmOrder 注文を確定（Confirmedステータスに変更）
// この時点で法的拘束力が発生し、コンプライアンスエンジンが動作する
func (s *OrderService) ConfirmOrder(ctx context.Context, req *ConfirmOrderRequest) (*domain.Order, error) {
	// 1. 既存の注文を取得
	oldOrder, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 2. セキュリティチェック: テナントIDが一致しているか
	if oldOrder.TenantID != req.TenantID {
		return nil, fmt.Errorf("unauthorized: tenant_id mismatch")
	}

	// 3. ステータスチェック: Draftステータスのみ確定可能
	if oldOrder.Status != domain.OrderStatusDraft {
		return nil, fmt.Errorf("order status must be Draft to confirm, current status: %s", oldOrder.Status)
	}

	// 4. コンプライアンス要件の検証
	if oldOrder.Details == nil || oldOrder.Details.Description == "" {
		return nil, fmt.Errorf("order details description is required for compliance")
	}

	// 5. ステータスをConfirmedに変更
	newOrder := *oldOrder
	newOrder.Status = domain.OrderStatusConfirmed
	newOrder.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, &newOrder); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// 6. 監査ログ記録（非同期・エラー時も継続）
	if s.auditLogRepo != nil {
		var ctxData *auditLogContext = &auditLogContext{
			TenantID:      req.TenantID,
			UserID:        req.UserID,
			Action:        domain.AuditActionConfirm,
			ResourceType:  "order",
			ResourceID:    newOrder.ID,
			OldValue:      s.orderToJSON(oldOrder),
			NewValue:      s.orderToJSON(&newOrder),
			ChangedFields: []string{"status"},
			IPAddress:     req.IPAddress,
			UserAgent:     req.UserAgent,
		}
		s.recordAuditLog(ctxData)
	}

	// 7. アンバサダー成果報酬を確定（非同期・エラー時も継続）
	if s.ambassadorService != nil {
		go func() {
			if err := s.ambassadorService.ApproveCommission(context.Background(), newOrder.ID); err != nil {
				fmt.Printf("WARNING: Failed to approve commission for order: %v\n", err)
			}
		}()
	}

	// 注意: コンプライアンスエンジン（PDF生成）は、別のサービス（Cloud Function）で
	// 非同期に実行される想定。ここではステータス変更のみを行う。

	return &newOrder, nil
}

// GetOrder 注文を取得
func (s *OrderService) GetOrder(ctx context.Context, orderID, tenantID string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// セキュリティチェック: テナントIDが一致しているか
	if order.TenantID != tenantID {
		return nil, fmt.Errorf("unauthorized: tenant_id mismatch")
	}

	return order, nil
}

// ListOrders 注文一覧を取得
func (s *OrderService) ListOrders(ctx context.Context, tenantID string) ([]*domain.Order, error) {
	orders, err := s.orderRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, nil
}

// auditLogContext 監査ログ記録用のコンテキスト
type auditLogContext struct {
	TenantID      string
	UserID        string
	Action        domain.AuditAction
	ResourceType  string
	ResourceID    string
	OldValue      string
	NewValue      string
	ChangedFields []string
	IPAddress     string
	UserAgent     string
}

// recordAuditLog 監査ログを記録（非同期、エラー時も継続）
func (s *OrderService) recordAuditLog(ctxData *auditLogContext) {
	go func() {
		// バックグラウンドで監査ログを記録
		// エラーが発生してもビジネスロジックには影響しない
		auditLog := domain.NewAuditLog(
			ctxData.TenantID,
			ctxData.UserID,
			ctxData.Action,
			ctxData.ResourceType,
			ctxData.ResourceID,
		)

		auditLog.OldValue = ctxData.OldValue
		auditLog.NewValue = ctxData.NewValue
		auditLog.ChangedFields = ctxData.ChangedFields
		auditLog.IPAddress = ctxData.IPAddress
		auditLog.UserAgent = ctxData.UserAgent

		if err := s.auditLogRepo.Create(context.Background(), auditLog); err != nil {
			// エラーログのみ記録（ビジネスロジックには影響しない）
			fmt.Printf("WARNING: Failed to record audit log: %v\n", err)
		}
	}()
}

// orderToJSON 注文オブジェクトをJSON文字列に変換（監査ログ用）
func (s *OrderService) orderToJSON(order *domain.Order) string {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal order: %v"}`, err)
	}
	return string(data)
}
