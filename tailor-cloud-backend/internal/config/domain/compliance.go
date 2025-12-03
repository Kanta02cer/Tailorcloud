package domain

import (
	"fmt"
	"time"
)

// ComplianceDocument コンプライアンス文書モデル
// 下請法第3条に基づく発注書の履歴管理
type ComplianceDocument struct {
	ID                string    `json:"id" db:"id"`
	OrderID           string    `json:"order_id" db:"order_id"`
	DocumentType      DocumentType `json:"document_type" db:"document_type"`
	ParentDocumentID  *string   `json:"parent_document_id" db:"parent_document_id"` // 修正元の文書ID
	PDFURL            string    `json:"pdf_url" db:"pdf_url"`
	PDFHash           string    `json:"pdf_hash" db:"pdf_hash"`
	GeneratedAt       time.Time `json:"generated_at" db:"generated_at"`
	GeneratedBy       string    `json:"generated_by" db:"generated_by"`
	AmendmentReason   *string   `json:"amendment_reason" db:"amendment_reason"` // 修正理由
	Version           int       `json:"version" db:"version"`                   // バージョン番号
	TenantID          string    `json:"tenant_id" db:"tenant_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// DocumentType 文書タイプ
type DocumentType string

const (
	DocumentTypeInitial    DocumentType = "INITIAL"    // 初回発注書
	DocumentTypeAmendment  DocumentType = "AMENDMENT"  // 修正発注書
)

// IsInitial 初回発注書かどうか
func (d *ComplianceDocument) IsInitial() bool {
	return d.DocumentType == DocumentTypeInitial
}

// IsAmendment 修正発注書かどうか
func (d *ComplianceDocument) IsAmendment() bool {
	return d.DocumentType == DocumentTypeAmendment
}

// HasParent 親文書があるかどうか
func (d *ComplianceDocument) HasParent() bool {
	return d.ParentDocumentID != nil && *d.ParentDocumentID != ""
}

// ComplianceRequirement コンプライアンス要件
// 下請法・フリーランス保護法で定められた必須項目
type ComplianceRequirement struct {
	// 委託をする者の氏名（発注者側の法人名）
	PrincipalName string `json:"principal_name"`
	
	// 給付の内容（業務内容の詳細）
	ServiceDescription string `json:"service_description"`
	
	// 報酬の額（税抜）
	RewardAmount int64 `json:"reward_amount"`
	
	// 支払期日（下請法60日ルール準拠）
	PaymentDueDate time.Time `json:"payment_due_date"`
	
	// 納期
	DeliveryDate time.Time `json:"delivery_date"`
}

// ValidateComplianceRequirement コンプライアンス要件の検証
func (cr *ComplianceRequirement) Validate() error {
	if cr.PrincipalName == "" {
		return fmt.Errorf("委託をする者の氏名が未指定です")
	}
	if cr.ServiceDescription == "" {
		return fmt.Errorf("給付の内容が未指定です")
	}
	if cr.RewardAmount <= 0 {
		return fmt.Errorf("報酬の額が無効です（0円以下）")
	}
	if cr.PaymentDueDate.IsZero() {
		return fmt.Errorf("支払期日が未指定です")
	}
	if cr.DeliveryDate.IsZero() {
		return fmt.Errorf("納期が未指定です")
	}
	
	// 下請法60日ルール: 納期から支払期日まで60日以上必要
	daysBetween := int(cr.PaymentDueDate.Sub(cr.DeliveryDate).Hours() / 24)
	if daysBetween < 60 {
		return fmt.Errorf("下請法60日ルール違反: 納期から支払期日まで60日以上必要（現在%d日）", daysBetween)
	}
	
	return nil
}

// BuildComplianceRequirementFromOrder 注文情報からコンプライアンス要件を構築
func BuildComplianceRequirementFromOrder(order *Order, tenant *Tenant, details *OrderDetails) *ComplianceRequirement {
	return &ComplianceRequirement{
		PrincipalName:      tenant.LegalName,
		ServiceDescription: details.Description,
		RewardAmount:       order.TotalAmount,
		PaymentDueDate:     order.PaymentDueDate,
		DeliveryDate:       order.DeliveryDate,
	}
}

