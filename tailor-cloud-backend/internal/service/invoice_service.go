package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// InvoiceService 請求書PDF生成サービス
// インボイス制度対応: 適格請求書（インボイス）のPDF生成
type InvoiceService struct {
	orderRepo    repository.OrderRepository
	tenantRepo   repository.TenantRepository
	customerRepo repository.CustomerRepository
	storageService StorageService
	bucketName   string
	taxService   *TaxCalculationService
}

// NewInvoiceService InvoiceServiceのコンストラクタ
func NewInvoiceService(
	orderRepo repository.OrderRepository,
	tenantRepo repository.TenantRepository,
	customerRepo repository.CustomerRepository,
	storageService StorageService,
	bucketName string,
	taxService *TaxCalculationService,
) *InvoiceService {
	return &InvoiceService{
		orderRepo:      orderRepo,
		tenantRepo:     tenantRepo,
		customerRepo:   customerRepo,
		storageService: storageService,
		bucketName:     bucketName,
		taxService:     taxService,
	}
}

// InvoiceRequest 請求書生成リクエスト
type InvoiceRequest struct {
	OrderID string
}

// InvoiceResponse 請求書生成レスポンス
type InvoiceResponse struct {
	OrderID      string
	InvoiceURL   string
	InvoiceHash  string
	IssuedAt     time.Time
	TaxAmount    int64
	TaxRate      domain.TaxRate
	TotalAmount  int64
}

// GenerateInvoice 適格請求書（インボイス）を生成
func (s *InvoiceService) GenerateInvoice(ctx context.Context, req *InvoiceRequest) (*InvoiceResponse, error) {
	// 注文情報を取得
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// テナント情報を取得（T番号、法人名、住所を取得）
	tenant, err := s.tenantRepo.GetByID(ctx, order.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// 顧客情報を取得
	var customer *domain.Customer
	if order.CustomerID != "" {
		customer, err = s.customerRepo.GetByID(ctx, order.CustomerID, order.TenantID)
		if err != nil {
			// 顧客情報が取得できない場合は警告のみ（請求書は生成可能）
			// log.Printf("WARNING: failed to get customer: %v", err)
		}
	}

	// 消費税額を計算（まだ計算されていない場合）
	var taxAmount int64
	var taxRate domain.TaxRate
	if order.TaxAmount == 0 {
		taxResp, err := s.taxService.CalculateTaxForOrder(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax: %w", err)
		}
		taxAmount = taxResp.TaxAmount
		taxRate = taxResp.TaxRate
	} else {
		taxAmount = order.TaxAmount
		taxRate = order.TaxRate
		if taxRate == 0 {
			taxRate = domain.TaxRateStandard
		}
	}

	// 税抜金額を取得
	taxExcludedAmount := order.TotalAmount
	if order.TaxExcludedAmount != nil {
		taxExcludedAmount = *order.TaxExcludedAmount
	}

	// PDFを生成
	pdfBytes, err := s.generateInvoicePDF(order, tenant, customer, taxExcludedAmount, taxAmount, taxRate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice PDF: %w", err)
	}

	// PDFのハッシュ値を計算（改ざん防止）
	hash := sha256.Sum256(pdfBytes)
	hashHex := hex.EncodeToString(hash[:])

	// Cloud Storageにアップロード
	objectPath := fmt.Sprintf("invoices/%s/invoice_%s_%s.pdf",
		order.TenantID,
		order.ID,
		time.Now().Format("20060102_150405"))

	docURL, err := s.storageService.UploadPDF(ctx, s.bucketName, objectPath, pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to upload invoice PDF: %w", err)
	}

	// 注文に請求書発行日時を記録（更新はOrderService経由で行う必要がある）
	issuedAt := time.Now()

	return &InvoiceResponse{
		OrderID:     order.ID,
		InvoiceURL:  docURL,
		InvoiceHash: hashHex,
		IssuedAt:    issuedAt,
		TaxAmount:   taxAmount,
		TaxRate:     taxRate,
		TotalAmount: taxExcludedAmount + taxAmount,
	}, nil
}

// generateInvoicePDF 適格請求書PDFを生成
func (s *InvoiceService) generateInvoicePDF(
	order *domain.Order,
	tenant *domain.Tenant,
	customer *domain.Customer,
	taxExcludedAmount int64,
	taxAmount int64,
	taxRate domain.TaxRate,
) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// 日本語フォント対応（注意: フォントファイルが必要）
	// 実際の本番環境では、Noto Sans JPなどの日本語フォントを埋め込む必要がある
	// ここでは基本のフォントを使用（後で日本語フォント対応が必要）
	
	// タイトル
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "INVOICE / REQUEST FOR PAYMENT", "", 0, "C", false, 0, "")
	pdf.Ln(15)

	// 発行事項
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Invoice No.: %s", order.ID), "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(0, 6, fmt.Sprintf("Issue Date: %s", time.Now().Format("2006-01-02")), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// 発行元（テナント）情報
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "FROM (Bill From)", "", 0, "L", false, 0, "")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	if tenant.LegalName != "" {
		pdf.CellFormat(0, 6, tenant.LegalName, "", 0, "L", false, 0, "")
		pdf.Ln(6)
	}
	if tenant.Address != "" {
		pdf.CellFormat(0, 6, tenant.Address, "", 0, "L", false, 0, "")
		pdf.Ln(6)
	}
	// インボイス登録番号（T番号）
	if tenant.InvoiceRegistrationNo != "" {
		pdf.CellFormat(0, 6, fmt.Sprintf("Invoice Registration No. (T-Number): %s", tenant.InvoiceRegistrationNo), "", 0, "L", false, 0, "")
		pdf.Ln(6)
	}
	pdf.Ln(5)

	// 宛先（顧客）情報
	if customer != nil {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 7, "TO (Bill To)", "", 0, "L", false, 0, "")
		pdf.Ln(7)
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, customer.Name, "", 0, "L", false, 0, "")
		pdf.Ln(6)
		if customer.Email != "" {
			pdf.CellFormat(0, 6, fmt.Sprintf("Email: %s", customer.Email), "", 0, "L", false, 0, "")
			pdf.Ln(6)
		}
		pdf.Ln(5)
	}

	// 明細
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 7, "DETAILS", "", 0, "L", false, 0, "")
	pdf.Ln(7)

	// テーブルヘッダー
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(30, 7, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 7, "Item", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Tax Excluded", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, "Tax Amount", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, "Total (Inc. Tax)", "1", 0, "R", true, 0, "")
	pdf.Ln(7)

	// 明細行
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	
	description := "Bespoke Suit Order"
	if order.Details != nil && order.Details.Description != "" {
		description = order.Details.Description
	}
	
	pdf.CellFormat(30, 7, "Order", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 7, description, "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s", formatCurrency(taxExcludedAmount)), "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s (%s)", formatCurrency(taxAmount), domain.FormatTaxRate(taxRate)), "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s", formatCurrency(taxExcludedAmount+taxAmount)), "1", 0, "R", false, 0, "")
	pdf.Ln(7)

	// 合計行
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(70, 7, "TOTAL", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s", formatCurrency(taxExcludedAmount)), "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s", formatCurrency(taxAmount)), "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("¥%s", formatCurrency(taxExcludedAmount+taxAmount)), "1", 0, "R", true, 0, "")
	pdf.Ln(10)

	// 支払条件
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("Payment Due Date: %s", order.PaymentDueDate.Format("2006-01-02")), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	if !order.DeliveryDate.IsZero() {
		pdf.CellFormat(0, 5, fmt.Sprintf("Delivery Date: %s", order.DeliveryDate.Format("2006-01-02")), "", 0, "L", false, 0, "")
		pdf.Ln(5)
	}

	// フッター
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated by TailorCloud ERP System on %s", time.Now().Format("2006-01-02 15:04:05")), "", 0, "C", false, 0, "")

	// PDFをバイト配列に変換
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF bytes: %w", err)
	}

	return buf.Bytes(), nil
}

// formatCurrencyはcompliance_service.goで定義されている共通関数を使用

