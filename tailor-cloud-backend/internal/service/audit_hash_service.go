package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"tailor-cloud/backend/internal/config/domain"
)

// AuditHashService 監査ログのハッシュ計算サービス
// 改ざん検知用
type AuditHashService struct{}

// NewAuditHashService AuditHashServiceのコンストラクタ
func NewAuditHashService() *AuditHashService {
	return &AuditHashService{}
}

// CalculateLogHash 監査ログのハッシュ値を計算
// ログの改ざん検知に使用
func (s *AuditHashService) CalculateLogHash(log *domain.AuditLog) (string, error) {
	// ハッシュ計算に使用するフィールドをJSON形式にシリアライズ
	hashData := map[string]interface{}{
		"id":            log.ID,
		"tenant_id":     log.TenantID,
		"user_id":       log.UserID,
		"action":        log.Action,
		"resource_type": log.ResourceType,
		"resource_id":   log.ResourceID,
		"old_value":     log.OldValue,
		"new_value":     log.NewValue,
		"ip_address":    log.IPAddress,
		"created_at":    log.CreatedAt,
	}
	
	data, err := json.Marshal(hashData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal log data: %w", err)
	}
	
	// SHA-256ハッシュを計算
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// VerifyLogHash ログのハッシュ値を検証
func (s *AuditHashService) VerifyLogHash(log *domain.AuditLog) (bool, error) {
	if log.LogHash == "" {
		return false, fmt.Errorf("log hash is empty")
	}
	
	calculatedHash, err := s.CalculateLogHash(log)
	if err != nil {
		return false, fmt.Errorf("failed to calculate hash: %w", err)
	}
	
	return calculatedHash == log.LogHash, nil
}

