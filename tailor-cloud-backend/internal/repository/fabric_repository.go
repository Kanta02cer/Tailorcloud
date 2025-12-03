package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
	"tailor-cloud/backend/internal/config/domain"
)

// FabricRepository 生地リポジトリインターフェース
type FabricRepository interface {
	GetByID(ctx context.Context, fabricID string) (*domain.Fabric, error)
	GetAll(ctx context.Context, tenantID string, filters *FabricFilters) ([]*domain.Fabric, error)
	Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Fabric, error)
	UpdateStock(ctx context.Context, fabricID string, stockAmount float64) error
}

// FabricFilters 生地フィルター
type FabricFilters struct {
	Status []domain.StockStatus // 在庫ステータスでフィルター（Available, Limited, SoldOut）
	Search string               // 検索キーワード（生地名で検索）
}

// PostgreSQLFabricRepository PostgreSQLを使った生地リポジトリ実装
type PostgreSQLFabricRepository struct {
	db *sql.DB
}

// NewPostgreSQLFabricRepository PostgreSQLFabricRepositoryのコンストラクタ
func NewPostgreSQLFabricRepository(db *sql.DB) FabricRepository {
	return &PostgreSQLFabricRepository{
		db: db,
	}
}

// GetByID 生地IDで取得
func (r *PostgreSQLFabricRepository) GetByID(ctx context.Context, fabricID string) (*domain.Fabric, error) {
	query := `
		SELECT 
			id, supplier_id, name, stock_amount, price,
			image_url, minimum_order,
			created_at, updated_at
		FROM fabrics
		WHERE id = $1
	`
	
	var fabric domain.Fabric
	err := r.db.QueryRowContext(ctx, query, fabricID).Scan(
		&fabric.ID,
		&fabric.SupplierID,
		&fabric.Name,
		&fabric.StockAmount,
		&fabric.Price,
		&fabric.ImageURL,
		&fabric.MinimumOrder,
		&fabric.CreatedAt,
		&fabric.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fabric not found: %s", fabricID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fabric: %w", err)
	}
	
	// 在庫ステータスを計算
	fabric.CalculateStockStatus()
	
	return &fabric, nil
}

// GetAll 生地一覧を取得（フィルター対応）
func (r *PostgreSQLFabricRepository) GetAll(ctx context.Context, tenantID string, filters *FabricFilters) ([]*domain.Fabric, error) {
	// クエリ構築
	query := `
		SELECT 
			id, supplier_id, name, stock_amount, price,
			image_url, minimum_order,
			created_at, updated_at
		FROM fabrics
		WHERE 1=1
	`
	
	args := []interface{}{}
	argIndex := 1
	
	// テナントIDフィルター（将来のマルチテナント対応）
	// 現時点では、fabricsテーブルにtenant_idが無いため、全件取得
	// Phase 2でtenant_idカラムを追加予定
	
	// 検索キーワードフィルター
	if filters != nil && filters.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}
	
	// 在庫ステータスフィルター
	if filters != nil && len(filters.Status) > 0 {
		// 在庫ステータスは計算フィールドなので、全件取得後にフィルター
		// または、在庫数量でフィルター（パフォーマンス最適化）
		statusConditions := []string{}
		for _, status := range filters.Status {
			switch status {
			case domain.StockStatusAvailable:
				statusConditions = append(statusConditions, fmt.Sprintf("stock_amount > 3.2"))
			case domain.StockStatusLimited:
				statusConditions = append(statusConditions, fmt.Sprintf("stock_amount > 0 AND stock_amount <= 3.2"))
			case domain.StockStatusSoldOut:
				statusConditions = append(statusConditions, "stock_amount = 0")
			}
		}
		if len(statusConditions) > 0 {
			query += " AND (" + strings.Join(statusConditions, " OR ") + ")"
		}
	}
	
	query += " ORDER BY name ASC"
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query fabrics: %w", err)
	}
	defer rows.Close()
	
	fabrics := make([]*domain.Fabric, 0)
	
	for rows.Next() {
		var fabric domain.Fabric
		
		err := rows.Scan(
			&fabric.ID,
			&fabric.SupplierID,
			&fabric.Name,
			&fabric.StockAmount,
			&fabric.Price,
			&fabric.ImageURL,
			&fabric.MinimumOrder,
			&fabric.CreatedAt,
			&fabric.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan fabric: %w", err)
		}
		
		// 在庫ステータスを計算
		fabric.CalculateStockStatus()
		
		// ステータスフィルター（在庫数量ベースのフィルターが適用されていない場合のみ）
		if filters != nil && len(filters.Status) > 0 {
			matched := false
			for _, status := range filters.Status {
				if fabric.StockStatus == status {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		
		fabrics = append(fabrics, &fabric)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fabrics: %w", err)
	}
	
	return fabrics, nil
}

// Search 生地名で検索
func (r *PostgreSQLFabricRepository) Search(ctx context.Context, tenantID string, keyword string) ([]*domain.Fabric, error) {
	filters := &FabricFilters{
		Search: keyword,
	}
	return r.GetAll(ctx, tenantID, filters)
}

// UpdateStock 在庫数量を更新
func (r *PostgreSQLFabricRepository) UpdateStock(ctx context.Context, fabricID string, stockAmount float64) error {
	query := `
		UPDATE fabrics SET
			stock_amount = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.ExecContext(ctx, query, fabricID, stockAmount)
	if err != nil {
		return fmt.Errorf("failed to update fabric stock: %w", err)
	}
	
	return nil
}

