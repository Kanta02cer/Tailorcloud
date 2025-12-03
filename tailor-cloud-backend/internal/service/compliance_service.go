package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf/v2"
	"tailor-cloud/backend/internal/config/domain"
)

// ComplianceService コンプライアンスサービス
// PDF生成とコンプライアンス検証を担当
type ComplianceService struct {
	storageService           StorageService
	bucketName               string // Cloud Storageバケット名
	jpFontHelper             *JPFontHelper // 日本語フォントヘルパー
	complianceDocRepo        ComplianceDocumentRepository // コンプライアンス文書リポジトリ
}

// ComplianceDocumentRepository コンプライアンス文書リポジトリインターフェース
type ComplianceDocumentRepository interface {
	Create(ctx context.Context, doc *domain.ComplianceDocument) error
	GetByID(ctx context.Context, docID string, tenantID string) (*domain.ComplianceDocument, error)
	GetByOrderID(ctx context.Context, orderID string, tenantID string) ([]*domain.ComplianceDocument, error)
	GetLatestByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error)
	GetInitialByOrderID(ctx context.Context, orderID string, tenantID string) (*domain.ComplianceDocument, error)
	GetVersionByOrderID(ctx context.Context, orderID string, tenantID string) (int, error)
}

// NewComplianceService ComplianceServiceのコンストラクタ
func NewComplianceService(storageService StorageService, bucketName string, complianceDocRepo ComplianceDocumentRepository) *ComplianceService {
	fontDir := GetFontDir()
	jpFontHelper := NewJPFontHelper(fontDir)
	
	return &ComplianceService{
		storageService:    storageService,
		bucketName:        bucketName,
		jpFontHelper:      jpFontHelper,
		complianceDocRepo: complianceDocRepo,
	}
}

// GenerateComplianceDocumentRequest PDF生成リクエスト
type GenerateComplianceDocumentRequest struct {
	Order          *domain.Order
	Tenant         *domain.Tenant
	Requirement    *domain.ComplianceRequirement
}

// GenerateComplianceDocumentResponse PDF生成レスポンス
type GenerateComplianceDocumentResponse struct {
	DocURL       string // Cloud Storage上のPDFのURL
	DocHash      string // SHA256ハッシュ値（改ざん防止用）
	DocumentID   string // コンプライアンス文書ID
	DocumentType domain.DocumentType // 文書タイプ（INITIAL or AMENDMENT）
	Version      int    // バージョン番号
}

// GenerateAmendmentDocumentRequest 修正発注書生成リクエスト
type GenerateAmendmentDocumentRequest struct {
	OrderID          string
	TenantID         string
	GeneratedBy      string // 発行者ユーザーID
	AmendmentReason  string // 修正理由
}

// GenerateAmendmentDocumentResponse 修正発注書生成レスポンス
type GenerateAmendmentDocumentResponse struct {
	DocumentID       string
	ParentDocumentID string
	DocURL           string
	DocHash          string
	Version          int
	AmendmentReason  string
}

// GenerateComplianceDocument コンプライアンスドキュメント（PDF）を生成
// 下請法・フリーランス保護法に準拠した発注書PDFを生成
func (s *ComplianceService) GenerateComplianceDocument(ctx context.Context, req *GenerateComplianceDocumentRequest) (*GenerateComplianceDocumentResponse, error) {
	// 1. コンプライアンス要件の検証
	if err := req.Requirement.Validate(); err != nil {
		return nil, fmt.Errorf("compliance requirement validation failed: %w", err)
	}
	
	// 2. PDF生成
	pdfBytes, err := s.generatePDF(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	// 3. PDFのハッシュ値を計算（改ざん防止）
	hash := sha256.Sum256(pdfBytes)
	hashHex := hex.EncodeToString(hash[:])
	
	// 4. Cloud Storageにアップロード
	objectPath := fmt.Sprintf("compliance-docs/%s/%s.pdf", req.Order.TenantID, req.Order.ID)
	var docURL string
	
	if s.storageService != nil && s.bucketName != "" {
		uploadedURL, err := s.storageService.UploadPDF(ctx, s.bucketName, objectPath, pdfBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to upload PDF to Cloud Storage: %w", err)
		}
		docURL = uploadedURL
	} else {
		// Storage Serviceが設定されていない場合はローカルパスのみ
		docURL = fmt.Sprintf("gs://%s/%s", s.bucketName, objectPath)
	}
	
	// 5. コンプライアンス文書レコードを作成（履歴管理）
	// 既存の文書があるかチェック
	existingDocs, _ := s.complianceDocRepo.GetByOrderID(ctx, req.Order.ID, req.Order.TenantID)
	var documentType domain.DocumentType
	var version int
	var parentDocID *string
	
	if len(existingDocs) == 0 {
		// 初回発注書
		documentType = domain.DocumentTypeInitial
		version = 1
		parentDocID = nil
	} else {
		// 修正発注書
		documentType = domain.DocumentTypeAmendment
		latestVersion, _ := s.complianceDocRepo.GetVersionByOrderID(ctx, req.Order.ID, req.Order.TenantID)
		version = latestVersion + 1
		
		// 最新の文書を親として設定
		latestDoc, _ := s.complianceDocRepo.GetLatestByOrderID(ctx, req.Order.ID, req.Order.TenantID)
		if latestDoc != nil {
			parentDocID = &latestDoc.ID
		}
	}
	
	// コンプライアンス文書を作成
	complianceDoc := &domain.ComplianceDocument{
		ID:               uuid.New().String(),
		OrderID:          req.Order.ID,
		DocumentType:     documentType,
		ParentDocumentID: parentDocID,
		PDFURL:           docURL,
		PDFHash:          hashHex,
		GeneratedAt:      time.Now(),
		GeneratedBy:      req.Order.CreatedBy,
		Version:          version,
		TenantID:         req.Order.TenantID,
	}
	
	// リポジトリに保存
	if s.complianceDocRepo != nil {
		if err := s.complianceDocRepo.Create(ctx, complianceDoc); err != nil {
			// エラーはログに記録するが、PDF生成自体は成功とみなす
			fmt.Printf("WARNING: Failed to save compliance document record: %v\n", err)
		}
	}
	
	return &GenerateComplianceDocumentResponse{
		DocURL:       docURL,
		DocHash:      hashHex,
		DocumentID:   complianceDoc.ID,
		DocumentType: documentType,
		Version:      version,
	}, nil
}

// generatePDF PDFを生成（下請法・フリーランス保護法準拠）
func (s *ComplianceService) generatePDF(req *GenerateComplianceDocumentRequest) ([]byte, error) {
	// PDF初期化（A4サイズ、縦向き）
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("発注書（下請法第3条）", false)
	pdf.SetAuthor("TailorCloud", false)
	pdf.AddPage()
	
	// 日本語フォントを登録
	if err := s.jpFontHelper.RegisterJPFonts(pdf); err != nil {
		// フォント登録に失敗した場合は警告のみ（Arialを使用）
		fmt.Printf("WARNING: Failed to register Japanese fonts: %v\n", err)
	}
	
	// タイトル（日本語）
	s.jpFontHelper.SetJPFont(pdf, "B", 16)
	pdf.CellFormat(190, 10, "発注書（下請法第3条）", "", 1, "C", false, 0, "")
	pdf.Ln(5)
	
	// サブタイトル（英語）
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(190, 6, "ORDER DOCUMENT (Subcontract Act Article 3)", "", 1, "C", false, 0, "")
	pdf.Ln(10)
	
	// 下請法第3条に基づく記載事項
	s.jpFontHelper.SetJPFont(pdf, "B", 12)
	pdf.CellFormat(190, 8, "下請法第3条に基づく記載事項", "", 1, "L", false, 0, "")
	pdf.Ln(5)
	
	// 1. 委託をする者の氏名及び住所
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "委託をする者の氏名", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	pdf.CellFormat(130, 8, req.Requirement.PrincipalName, "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// 住所
	if req.Tenant != nil && req.Tenant.Address != "" {
		s.jpFontHelper.SetJPFont(pdf, "", 11)
		pdf.CellFormat(60, 8, "住所", "", 0, "L", false, 0, "")
		s.jpFontHelper.SetJPFont(pdf, "B", 11)
		pdf.CellFormat(130, 8, req.Tenant.Address, "", 1, "L", false, 0, "")
		pdf.Ln(3)
	}
	
	// 2. 給付の内容
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "給付の内容", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	pdf.CellFormat(130, 8, req.Requirement.ServiceDescription, "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// 3. 報酬の額
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "報酬の額（税抜）", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	amountStr := fmt.Sprintf("¥%s", formatCurrency(req.Requirement.RewardAmount))
	pdf.CellFormat(130, 8, amountStr, "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// 消費税（10%と仮定）
	taxAmount := req.Requirement.RewardAmount * 10 / 100
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "消費税（10%）", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	taxStr := fmt.Sprintf("¥%s", formatCurrency(taxAmount))
	pdf.CellFormat(130, 8, taxStr, "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// 合計金額
	totalAmount := req.Requirement.RewardAmount + taxAmount
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	pdf.CellFormat(60, 8, "合計金額（税込）", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 12)
	totalStr := fmt.Sprintf("¥%s", formatCurrency(totalAmount))
	pdf.CellFormat(130, 8, totalStr, "", 1, "L", false, 0, "")
	pdf.Ln(5)
	
	// 4. 納期
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "納期", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	deliveryDateStr := req.Requirement.DeliveryDate.Format("2006年01月02日")
	pdf.CellFormat(130, 8, deliveryDateStr, "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	// 5. 支払期日
	s.jpFontHelper.SetJPFont(pdf, "", 11)
	pdf.CellFormat(60, 8, "支払期日", "", 0, "L", false, 0, "")
	s.jpFontHelper.SetJPFont(pdf, "B", 11)
	paymentDateStr := req.Requirement.PaymentDueDate.Format("2006年01月02日")
	pdf.CellFormat(130, 8, paymentDateStr, "", 1, "L", false, 0, "")
	pdf.Ln(5)
	
	// 注文番号（Order.IDの末尾8文字）
	pdf.SetFont("Arial", "", 9)
	orderIDShort := req.Order.ID
	if len(orderIDShort) > 8 {
		orderIDShort = orderIDShort[len(orderIDShort)-8:]
	}
	pdf.CellFormat(190, 8, fmt.Sprintf("発注番号: %s", orderIDShort), "", 1, "R", false, 0, "")
	
	// 発行日時
	pdf.CellFormat(190, 8, fmt.Sprintf("発行日時: %s", time.Now().Format("2006年01月02日 15時04分05秒")), "", 1, "R", false, 0, "")
	
	// PDFをバイト配列に変換
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}
	
	return buf.Bytes(), nil
}

// formatCurrency 金額をカンマ区切りでフォーマット
func formatCurrency(amount int64) string {
	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return str
	}
	
	result := ""
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(str[i])
	}
	return result
}

// ValidateComplianceRequirement コンプライアンス要件を検証
func (s *ComplianceService) ValidateComplianceRequirement(req *domain.ComplianceRequirement) error {
	return req.Validate()
}

// CalculatePaymentDueDate 納期から支払期日を計算（下請法60日ルール）
func CalculatePaymentDueDate(deliveryDate time.Time) time.Time {
	return deliveryDate.AddDate(0, 0, 60)
}

// IsPaymentDueDateCompliant 支払期日が下請法に準拠しているかチェック
func IsPaymentDueDateCompliant(deliveryDate, paymentDueDate time.Time) bool {
	daysBetween := int(paymentDueDate.Sub(deliveryDate).Hours() / 24)
	return daysBetween >= 60
}

// GenerateAmendmentDocument 修正発注書を生成
// 元の発注書を上書きせず、新しいPDFを作成して履歴として保存
func (s *ComplianceService) GenerateAmendmentDocument(
	ctx context.Context,
	req *GenerateAmendmentDocumentRequest,
	order *domain.Order,
	tenant *domain.Tenant,
) (*GenerateAmendmentDocumentResponse, error) {
	// 1. 既存の最新発注書を取得（親文書として使用）
	latestDoc, err := s.complianceDocRepo.GetLatestByOrderID(ctx, req.OrderID, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest compliance document: %w", err)
	}
	
	// 2. コンプライアンス要件を構築（注文情報から）
	requirement := domain.BuildComplianceRequirementFromOrder(order, tenant, order.Details)
	if requirement == nil {
		return nil, fmt.Errorf("failed to build compliance requirement")
	}
	
	// 3. PDF生成
	pdfReq := &GenerateComplianceDocumentRequest{
		Order:       order,
		Tenant:      tenant,
		Requirement: requirement,
	}
	
	pdfBytes, err := s.generatePDF(pdfReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	// 4. PDFのハッシュ値を計算
	hash := sha256.Sum256(pdfBytes)
	hashHex := hex.EncodeToString(hash[:])
	
	// 5. Cloud Storageにアップロード
	version, _ := s.complianceDocRepo.GetVersionByOrderID(ctx, req.OrderID, req.TenantID)
	objectPath := fmt.Sprintf("compliance-docs/%s/%s_amendment_v%d_%s.pdf",
		req.TenantID,
		req.OrderID,
		version+1,
		time.Now().Format("20060102_150405"))
	
	var docURL string
	if s.storageService != nil && s.bucketName != "" {
		uploadedURL, err := s.storageService.UploadPDF(ctx, s.bucketName, objectPath, pdfBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to upload PDF to Cloud Storage: %w", err)
		}
		docURL = uploadedURL
	} else {
		docURL = fmt.Sprintf("gs://%s/%s", s.bucketName, objectPath)
	}
	
	// 6. コンプライアンス文書レコードを作成（修正発注書）
	complianceDoc := &domain.ComplianceDocument{
		ID:               uuid.New().String(),
		OrderID:          req.OrderID,
		DocumentType:     domain.DocumentTypeAmendment,
		ParentDocumentID: &latestDoc.ID,
		PDFURL:           docURL,
		PDFHash:          hashHex,
		GeneratedAt:      time.Now(),
		GeneratedBy:      req.GeneratedBy,
		AmendmentReason:  &req.AmendmentReason,
		Version:          version + 1,
		TenantID:         req.TenantID,
	}
	
	if err := s.complianceDocRepo.Create(ctx, complianceDoc); err != nil {
		return nil, fmt.Errorf("failed to create compliance document record: %w", err)
	}
	
	return &GenerateAmendmentDocumentResponse{
		DocumentID:       complianceDoc.ID,
		ParentDocumentID: latestDoc.ID,
		DocURL:           docURL,
		DocHash:          hashHex,
		Version:          complianceDoc.Version,
		AmendmentReason:  req.AmendmentReason,
	}, nil
}
