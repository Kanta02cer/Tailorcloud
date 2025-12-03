package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"tailor-cloud/backend/internal/config/domain"
)

// ComplianceDocumentRepository コンプライアンス文書リポジトリインターフェース
type ComplianceDocumentRepository interface {
	Create(ctx context.Context, doc *domain.ComplianceDocument) error
	GetByID(ctx context.Context, docID string, tenantID string) (*domain.ComplianceDocument, error)
	GetByOrderID(ctx context.Context, orderID string, tenantID string) ([]*domain.ComplianceDocument, error)
	GetLatestByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error)
	GetInitialByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error)
	GetVersionByOrderID(ctx context.Context, orderID string, tenantID string) (int, error)
}

// PostgreSQLComplianceDocumentRepository PostgreSQL実装
type PostgreSQLComplianceDocumentRepository struct {
	db *sql.DB
}

// NewPostgreSQLComplianceDocumentRepository PostgreSQLComplianceDocumentRepositoryのコンストラクタ
func NewPostgreSQLComplianceDocumentRepository(db *sql.DB) ComplianceDocumentRepository {
	return &PostgreSQLComplianceDocumentRepository{
		db: db,
	}
}

// Create コンプライアンス文書を作成
func (r *PostgreSQLComplianceDocumentRepository) Create(ctx context.Context, doc *domain.ComplianceDocument) error {
	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}
	
	now := time.Now()
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = now
	}
	if doc.UpdatedAt.IsZero() {
		doc.UpdatedAt = now
	}
	if doc.GeneratedAt.IsZero() {
		doc.GeneratedAt = now
	}
	
	query := `
		INSERT INTO compliance_documents (
			id, order_id, document_type, parent_document_id,
			pdf_url, pdf_hash, generated_at, generated_by,
			amendment_reason, version, tenant_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	
	var parentDocID interface{}
	if doc.ParentDocumentID != nil {
		parentDocID = *doc.ParentDocumentID
	}
	
	var amendmentReason interface{}
	if doc.AmendmentReason != nil {
		amendmentReason = *doc.AmendmentReason
	}
	
	_, err := r.db.ExecContext(ctx, query,
		doc.ID,
		doc.OrderID,
		string(doc.DocumentType),
		parentDocID,
		doc.PDFURL,
		doc.PDFHash,
		doc.GeneratedAt,
		doc.GeneratedBy,
		amendmentReason,
		doc.Version,
		doc.TenantID,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create compliance document: %w", err)
	}
	
	return nil
}

// GetByID コンプライアンス文書をIDで取得
func (r *PostgreSQLComplianceDocumentRepository) GetByID(ctx context.Context, docID string, tenantID string) (*domain.ComplianceDocument, error) {
	query := `
		SELECT 
			id, order_id, document_type, parent_document_id,
			pdf_url, pdf_hash, generated_at, generated_by,
			amendment_reason, version, tenant_id, created_at, updated_at
		FROM compliance_documents
		WHERE id = $1 AND tenant_id = $2
	`
	
	var doc domain.ComplianceDocument
	var documentTypeStr string
	var parentDocID, amendmentReason sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, docID, tenantID).Scan(
		&doc.ID,
		&doc.OrderID,
		&documentTypeStr,
		&parentDocID,
		&doc.PDFURL,
		&doc.PDFHash,
		&doc.GeneratedAt,
		&doc.GeneratedBy,
		&amendmentReason,
		&doc.Version,
		&doc.TenantID,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("compliance document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get compliance document: %w", err)
	}
	
	doc.DocumentType = domain.DocumentType(documentTypeStr)
	if parentDocID.Valid {
		doc.ParentDocumentID = &parentDocID.String
	}
	if amendmentReason.Valid {
		doc.AmendmentReason = &amendmentReason.String
	}
	
	return &doc, nil
}

// GetByOrderID 注文IDでコンプライアンス文書一覧を取得
func (r *PostgreSQLComplianceDocumentRepository) GetByOrderID(ctx context.Context, orderID string, tenantID string) ([]*domain.ComplianceDocument, error) {
	query := `
		SELECT 
			id, order_id, document_type, parent_document_id,
			pdf_url, pdf_hash, generated_at, generated_by,
			amendment_reason, version, tenant_id, created_at, updated_at
		FROM compliance_documents
		WHERE order_id = $1 AND tenant_id = $2
		ORDER BY version ASC, generated_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, orderID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list compliance documents: %w", err)
	}
	defer rows.Close()
	
	var documents []*domain.ComplianceDocument
	for rows.Next() {
		var doc domain.ComplianceDocument
		var documentTypeStr string
		var parentDocID, amendmentReason sql.NullString
		
		err := rows.Scan(
			&doc.ID,
			&doc.OrderID,
			&documentTypeStr,
			&parentDocID,
			&doc.PDFURL,
			&doc.PDFHash,
			&doc.GeneratedAt,
			&doc.GeneratedBy,
			&amendmentReason,
			&doc.Version,
			&doc.TenantID,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan compliance document: %w", err)
		}
		
		doc.DocumentType = domain.DocumentType(documentTypeStr)
		if parentDocID.Valid {
			doc.ParentDocumentID = &parentDocID.String
		}
		if amendmentReason.Valid {
			doc.AmendmentReason = &amendmentReason.String
		}
		
		documents = append(documents, &doc)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate compliance documents: %w", err)
	}
	
	return documents, nil
}

// GetLatestByOrderID 注文IDで最新のコンプライアンス文書を取得
func (r *PostgreSQLComplianceDocumentRepository) GetLatestByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error) {
	query := `
		SELECT 
			id, order_id, document_type, parent_document_id,
			pdf_url, pdf_hash, generated_at, generated_by,
			amendment_reason, version, tenant_id, created_at, updated_at
		FROM compliance_documents
		WHERE order_id = $1 AND tenant_id = $2
		ORDER BY version DESC, generated_at DESC
		LIMIT 1
	`
	
	var doc domain.ComplianceDocument
	var documentTypeStr string
	var parentDocID, amendmentReason sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, orderID, tenantID).Scan(
		&doc.ID,
		&doc.OrderID,
		&documentTypeStr,
		&parentDocID,
		&doc.PDFURL,
		&doc.PDFHash,
		&doc.GeneratedAt,
		&doc.GeneratedBy,
		&amendmentReason,
		&doc.Version,
		&doc.TenantID,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("compliance document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest compliance document: %w", err)
	}
	
	doc.DocumentType = domain.DocumentType(documentTypeStr)
	if parentDocID.Valid {
		doc.ParentDocumentID = &parentDocID.String
	}
	if amendmentReason.Valid {
		doc.AmendmentReason = &amendmentReason.String
	}
	
	return &doc, nil
}

// GetInitialByOrderID 注文IDで初回発注書を取得
func (r *PostgreSQLComplianceDocumentRepository) GetInitialByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error) {
	query := `
		SELECT 
			id, order_id, document_type, parent_document_id,
			pdf_url, pdf_hash, generated_at, generated_by,
			amendment_reason, version, tenant_id, created_at, updated_at
		FROM compliance_documents
		WHERE order_id = $1 AND tenant_id = $2 AND document_type = 'INITIAL'
		ORDER BY version ASC
		LIMIT 1
	`
	
	var doc domain.ComplianceDocument
	var documentTypeStr string
	var parentDocID, amendmentReason sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, orderID, tenantID).Scan(
		&doc.ID,
		&doc.OrderID,
		&documentTypeStr,
		&parentDocID,
		&doc.PDFURL,
		&doc.PDFHash,
		&doc.GeneratedAt,
		&doc.GeneratedBy,
		&amendmentReason,
		&doc.Version,
		&doc.TenantID,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("initial compliance document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get initial compliance document: %w", err)
	}
	
	doc.DocumentType = domain.DocumentType(documentTypeStr)
	if parentDocID.Valid {
		doc.ParentDocumentID = &parentDocID.String
	}
	if amendmentReason.Valid {
		doc.AmendmentReason = &amendmentReason.String
	}
	
	return &doc, nil
}

// GetVersionByOrderID 注文IDで最新のバージョン番号を取得
func (r *PostgreSQLComplianceDocumentRepository) GetVersionByOrderID(ctx context.Context, orderID string, tenantID string) (int, error) {
	query := `
		SELECT COALESCE(MAX(version), 0)
		FROM compliance_documents
		WHERE order_id = $1 AND tenant_id = $2
	`
	
	var version int
	err := r.db.QueryRowContext(ctx, query, orderID, tenantID).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to get version: %w", err)
	}
	
	return version, nil
}

