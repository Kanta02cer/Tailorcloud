package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"tailor-cloud/backend/internal/config/domain"
)

// PostgreSQLOrderRepository PostgreSQLを使った注文リポジトリ実装
// 仕様書に基づき、注文データはPrimary DB（PostgreSQL）に保存
type PostgreSQLOrderRepository struct {
	db *sql.DB
}

// NewPostgreSQLOrderRepository PostgreSQLOrderRepositoryのコンストラクタ
func NewPostgreSQLOrderRepository(db *sql.DB) OrderRepository {
	return &PostgreSQLOrderRepository{
		db: db,
	}
}

// Create 注文を作成
func (r *PostgreSQLOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (
			id, tenant_id, customer_id, fabric_id, status,
			compliance_doc_url, compliance_doc_hash,
			total_amount, payment_due_date, delivery_date,
			measurement_data, adjustments, description,
			created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	
	// OrderDetailsのJSONデータを準備
	var measurementDataJSON, adjustmentsJSON []byte
	var err error
	
	if order.Details != nil {
		if order.Details.MeasurementData != nil {
			measurementDataJSON = order.Details.MeasurementData
		}
		if order.Details.Adjustments != nil {
			adjustmentsJSON = order.Details.Adjustments
		}
	}
	
	description := ""
	if order.Details != nil {
		description = order.Details.Description
	}
	
	_, err = r.db.ExecContext(ctx, query,
		order.ID,
		order.TenantID,
		order.CustomerID,
		order.FabricID,
		string(order.Status),
		order.ComplianceDocURL,
		order.ComplianceDocHash,
		order.TotalAmount,
		order.PaymentDueDate,
		order.DeliveryDate,
		measurementDataJSON,
		adjustmentsJSON,
		description,
		order.CreatedAt,
		order.UpdatedAt,
		order.CreatedBy,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create order in postgresql: %w", err)
	}
	
	return nil
}

// GetByID 注文IDで取得
func (r *PostgreSQLOrderRepository) GetByID(ctx context.Context, orderID string) (*domain.Order, error) {
	query := `
		SELECT 
			id, tenant_id, customer_id, fabric_id, status,
			compliance_doc_url, compliance_doc_hash,
			total_amount, payment_due_date, delivery_date,
			measurement_data, adjustments, description,
			created_at, updated_at, created_by
		FROM orders
		WHERE id = $1
	`
	
	var order domain.Order
	var statusStr string
	var measurementDataJSON, adjustmentsJSON sql.NullString
	var description sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.TenantID,
		&order.CustomerID,
		&order.FabricID,
		&statusStr,
		&order.ComplianceDocURL,
		&order.ComplianceDocHash,
		&order.TotalAmount,
		&order.PaymentDueDate,
		&order.DeliveryDate,
		&measurementDataJSON,
		&adjustmentsJSON,
		&description,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CreatedBy,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order from postgresql: %w", err)
	}
	
	order.Status = domain.OrderStatus(statusStr)
	
	// OrderDetailsを構築
	if measurementDataJSON.Valid || adjustmentsJSON.Valid || description.Valid {
		order.Details = &domain.OrderDetails{}
		
		if measurementDataJSON.Valid {
			order.Details.MeasurementData = json.RawMessage(measurementDataJSON.String)
		}
		if adjustmentsJSON.Valid {
			order.Details.Adjustments = json.RawMessage(adjustmentsJSON.String)
		}
		if description.Valid {
			order.Details.Description = description.String
		}
	}
	
	return &order, nil
}

// GetByTenantID テナントIDで注文一覧を取得（後方互換性のため残す）
func (r *PostgreSQLOrderRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*domain.Order, error) {
	// ページネーションなしの取得（全件取得）
	return r.GetByTenantIDWithPagination(ctx, tenantID, 1, 10000) // 実質全件
}

// GetByTenantIDWithPagination テナントIDで注文一覧を取得（ページネーション対応）
func (r *PostgreSQLOrderRepository) GetByTenantIDWithPagination(ctx context.Context, tenantID string, page, pageSize int) ([]*domain.Order, error) {
	// オフセットとリミットを計算
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
	
	query := `
		SELECT 
			id, tenant_id, customer_id, fabric_id, status,
			compliance_doc_url, compliance_doc_hash,
			total_amount, payment_due_date, delivery_date,
			measurement_data, adjustments, description,
			created_at, updated_at, created_by
		FROM orders
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, tenantID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders by tenant_id: %w", err)
	}
	defer rows.Close()
	
	orders := make([]*domain.Order, 0)
	
	for rows.Next() {
		var order domain.Order
		var statusStr string
		var measurementDataJSON, adjustmentsJSON sql.NullString
		var description sql.NullString
		
		err := rows.Scan(
			&order.ID,
			&order.TenantID,
			&order.CustomerID,
			&order.FabricID,
			&statusStr,
			&order.ComplianceDocURL,
			&order.ComplianceDocHash,
			&order.TotalAmount,
			&order.PaymentDueDate,
			&order.DeliveryDate,
			&measurementDataJSON,
			&adjustmentsJSON,
			&description,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CreatedBy,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		
		order.Status = domain.OrderStatus(statusStr)
		
		// OrderDetailsを構築
		if measurementDataJSON.Valid || adjustmentsJSON.Valid || description.Valid {
			order.Details = &domain.OrderDetails{}
			
			if measurementDataJSON.Valid {
				order.Details.MeasurementData = json.RawMessage(measurementDataJSON.String)
			}
			if adjustmentsJSON.Valid {
				order.Details.Adjustments = json.RawMessage(adjustmentsJSON.String)
			}
			if description.Valid {
				order.Details.Description = description.String
			}
		}
		
		orders = append(orders, &order)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}
	
	return orders, nil
}

// CountByTenantID テナントIDで注文数を取得（ページネーション用）
func (r *PostgreSQLOrderRepository) CountByTenantID(ctx context.Context, tenantID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM orders
		WHERE tenant_id = $1
	`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count orders by tenant_id: %w", err)
	}
	
	return count, nil
}

// Update 注文を更新
func (r *PostgreSQLOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	// マルチテナント分離: 更新時もtenant_idが一致しているか確認
	existingOrder, err := r.GetByID(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing order: %w", err)
	}
	
	// セキュリティチェック: テナントIDが一致しているか
	if existingOrder.TenantID != order.TenantID {
		return fmt.Errorf("unauthorized: tenant_id mismatch")
	}
	
	query := `
		UPDATE orders SET
			customer_id = $2,
			fabric_id = $3,
			status = $4,
			compliance_doc_url = $5,
			compliance_doc_hash = $6,
			total_amount = $7,
			payment_due_date = $8,
			delivery_date = $9,
			measurement_data = $10,
			adjustments = $11,
			description = $12,
			updated_at = $13
		WHERE id = $1 AND tenant_id = $14
	`
	
	// OrderDetailsのJSONデータを準備
	var measurementDataJSON, adjustmentsJSON []byte
	var description string
	
	if order.Details != nil {
		if order.Details.MeasurementData != nil {
			measurementDataJSON = order.Details.MeasurementData
		}
		if order.Details.Adjustments != nil {
			adjustmentsJSON = order.Details.Adjustments
		}
		description = order.Details.Description
	}
	
	order.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.FabricID,
		string(order.Status),
		order.ComplianceDocURL,
		order.ComplianceDocHash,
		order.TotalAmount,
		order.PaymentDueDate,
		order.DeliveryDate,
		measurementDataJSON,
		adjustmentsJSON,
		description,
		order.UpdatedAt,
		order.TenantID, // WHERE句でテナントIDを確認
	)
	
	if err != nil {
		return fmt.Errorf("failed to update order in postgresql: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("order not found or tenant_id mismatch")
	}
	
	return nil
}

// UpdateStatus 注文ステータスを更新
func (r *PostgreSQLOrderRepository) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	query := `
		UPDATE orders SET
			status = $2,
			updated_at = $3
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, orderID, string(status), time.Now())
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	return nil
}

