package domain

import (
	"fmt"
	"math"
)

// TaxRate 消費税率
type TaxRate float64

const (
	TaxRateStandard TaxRate = 0.10 // 標準税率（10%）
	TaxRateReduced  TaxRate = 0.08 // 軽減税率（8%）
)

// TaxRoundingMethod 端数処理方法
type TaxRoundingMethod string

const (
	TaxRoundingMethodHalfUp   TaxRoundingMethod = "HALF_UP"   // 四捨五入
	TaxRoundingMethodRoundDown TaxRoundingMethod = "ROUND_DOWN" // 切り捨て
	TaxRoundingMethodRoundUp   TaxRoundingMethod = "ROUND_UP"   // 切り上げ
)

// CalculateTax 消費税額を計算
// taxExcludedAmount: 税抜金額
// taxRate: 消費税率（0.10 = 10%, 0.08 = 8%）
// roundingMethod: 端数処理方法
func CalculateTax(taxExcludedAmount int64, taxRate TaxRate, roundingMethod TaxRoundingMethod) int64 {
	// 税額 = 税抜金額 × 税率
	taxAmount := float64(taxExcludedAmount) * float64(taxRate)
	
	// 端数処理
	switch roundingMethod {
	case TaxRoundingMethodHalfUp:
		// 四捨五入
		return int64(math.Round(taxAmount))
	case TaxRoundingMethodRoundDown:
		// 切り捨て
		return int64(math.Floor(taxAmount))
	case TaxRoundingMethodRoundUp:
		// 切り上げ
		return int64(math.Ceil(taxAmount))
	default:
		// デフォルトは四捨五入
		return int64(math.Round(taxAmount))
	}
}

// CalculateTaxIncludedAmount 税込金額を計算
// taxExcludedAmount: 税抜金額
// taxRate: 消費税率
// roundingMethod: 端数処理方法
func CalculateTaxIncludedAmount(taxExcludedAmount int64, taxRate TaxRate, roundingMethod TaxRoundingMethod) int64 {
	taxAmount := CalculateTax(taxExcludedAmount, taxRate, roundingMethod)
	return taxExcludedAmount + taxAmount
}

// ParseTaxRate 税率文字列をパース
func ParseTaxRate(rateStr string) (TaxRate, error) {
	switch rateStr {
	case "0.10", "10", "10%":
		return TaxRateStandard, nil
	case "0.08", "8", "8%":
		return TaxRateReduced, nil
	default:
		return 0, fmt.Errorf("invalid tax rate: %s (must be 0.10 or 0.08)", rateStr)
	}
}

// FormatTaxRate 税率を表示用文字列にフォーマット
func FormatTaxRate(taxRate TaxRate) string {
	if taxRate == TaxRateStandard {
		return "10%"
	} else if taxRate == TaxRateReduced {
		return "8%"
	}
	return fmt.Sprintf("%.0f%%", taxRate*100)
}

