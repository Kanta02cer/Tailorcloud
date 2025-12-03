package repository

import (
	"context"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// OrderRepository 注文リポジトリインターフェース
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, orderID string) (*domain.Order, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Order, error)
	GetByTenantIDWithPagination(ctx context.Context, tenantID string, page, pageSize int) ([]*domain.Order, error)
	CountByTenantID(ctx context.Context, tenantID string) (int, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error
}

// FirestoreOrderRepository Firestoreを使った注文リポジトリ実装
type FirestoreOrderRepository struct {
	client *firestore.Client
}

// NewFirestoreOrderRepository FirestoreOrderRepositoryのコンストラクタ
func NewFirestoreOrderRepository(client *firestore.Client) OrderRepository {
	return &FirestoreOrderRepository{
		client: client,
	}
}

// Create 注文を作成
func (r *FirestoreOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	docRef := r.client.Collection("orders").Doc(order.ID)

	// Firestoreに保存（マルチテナント分離: tenant_idでクエリ可能にする）
	_, err := docRef.Set(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create order in firestore: %w", err)
	}

	return nil
}

// GetByID 注文IDで取得
func (r *FirestoreOrderRepository) GetByID(ctx context.Context, orderID string) (*domain.Order, error) {
	docRef := r.client.Collection("orders").Doc(orderID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from firestore: %w", err)
	}

	var order domain.Order
	if err := docSnap.DataTo(&order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

// GetByTenantID テナントIDで注文一覧を取得（後方互換性のため残す）
func (r *FirestoreOrderRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Order, error) {
	return r.GetByTenantIDWithPagination(ctx, tenantID, 1, 10000) // 実質全件
}

// GetByTenantIDWithPagination テナントIDで注文一覧を取得（ページネーション対応）
func (r *FirestoreOrderRepository) GetByTenantIDWithPagination(ctx context.Context, tenantID string, page, pageSize int) ([]*domain.Order, error) {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	query := r.client.Collection("orders").
		Where("tenant_id", "==", tenantID).
		OrderBy("created_at", firestore.Desc).
		Offset(offset).
		Limit(pageSize)
	
	iter := query.Documents(ctx)
	defer iter.Stop()
	
	orders := make([]*domain.Order, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, fmt.Errorf("failed to iterate orders: %w", err)
		}
		
		var order domain.Order
		if err := doc.DataTo(&order); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %w", err)
		}
		
		orders = append(orders, &order)
	}
	
	return orders, nil
}

// CountByTenantID テナントIDで注文数を取得（ページネーション用）
func (r *FirestoreOrderRepository) CountByTenantID(ctx context.Context, tenantID string) (int, error) {
	query := r.client.Collection("orders").Where("tenant_id", "==", tenantID)
	
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to count orders: %w", err)
	}
	
	return len(docs), nil
}

// Update 注文を更新
func (r *FirestoreOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	docRef := r.client.Collection("orders").Doc(order.ID)

	// マルチテナント分離: 更新時もtenant_idが一致しているか確認
	existingDoc, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing order: %w", err)
	}

	var existingOrder domain.Order
	if err := existingDoc.DataTo(&existingOrder); err != nil {
		return fmt.Errorf("failed to unmarshal existing order: %w", err)
	}

	// セキュリティチェック: テナントIDが一致しているか
	if existingOrder.TenantID != order.TenantID {
		return fmt.Errorf("unauthorized: tenant_id mismatch")
	}

	// 更新実行
	_, err = docRef.Set(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order in firestore: %w", err)
	}

	return nil
}

// UpdateStatus 注文ステータスを更新
func (r *FirestoreOrderRepository) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	docRef := r.client.Collection("orders").Doc(orderID)

	// 部分更新（ステータスのみ）
	updates := []firestore.Update{
		{Path: "status", Value: string(status)},
		{Path: "updated_at", Value: firestore.ServerTimestamp},
	}

	_, err := docRef.Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
