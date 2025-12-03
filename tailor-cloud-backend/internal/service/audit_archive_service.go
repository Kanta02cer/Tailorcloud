package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/config/domain"
	"tailor-cloud/backend/internal/repository"
)

// AuditArchiveService 監査ログアーカイブサービス
// 1年以上のログをWORMストレージへ移行
type AuditArchiveService struct {
	auditLogRepo   repository.AuditLogRepository
	storageService StorageService
	bucketName     string
	archiveRepo    repository.AuditLogArchiveRepository
}

// NewAuditArchiveService AuditArchiveServiceのコンストラクタ
func NewAuditArchiveService(
	auditLogRepo repository.AuditLogRepository,
	storageService StorageService,
	bucketName string,
	archiveRepo repository.AuditLogArchiveRepository,
) *AuditArchiveService {
	return &AuditArchiveService{
		auditLogRepo:  auditLogRepo,
		storageService: storageService,
		bucketName:    bucketName,
		archiveRepo:   archiveRepo,
	}
}

// ArchiveOldLogs 古いログをアーカイブ
// retentionPeriod: 保持期間（例: 365日）
func (s *AuditArchiveService) ArchiveOldLogs(ctx context.Context, tenantID string, retentionPeriodDays int) error {
	// アーカイブ対象の期間を決定（現在日時からretentionPeriodDays日前まで）
	cutoffDate := time.Now().AddDate(0, 0, -retentionPeriodDays)
	
	// 古いログを取得（1ヶ月単位でアーカイブ）
	startDate := cutoffDate.AddDate(0, -1, 0) // 1ヶ月前から
	endDate := cutoffDate
	
	// テナントごとのログを取得
	logs, err := s.auditLogRepo.GetLogsByDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get old logs: %w", err)
	}
	
	if len(logs) == 0 {
		// アーカイブ対象のログがない
		return nil
	}
	
	// ログをJSON形式にシリアライズ
	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}
	
	// アーカイブファイルのハッシュ値を計算
	hash := sha256.Sum256(logsJSON)
	hashHex := hex.EncodeToString(hash[:])
	
	// Cloud Storageにアップロード（WORM設定）
	objectPath := fmt.Sprintf("audit-archives/%s/%s_%s.json",
		tenantID,
		startDate.Format("200601"),
		endDate.Format("200601"))
	
	archiveURL, err := s.storageService.UploadJSON(ctx, s.bucketName, objectPath, logsJSON)
	if err != nil {
		return fmt.Errorf("failed to upload archive to Cloud Storage: %w", err)
	}
	
	// アーカイブメタデータを作成
	archive := &domain.AuditLogArchive{
		ID:                fmt.Sprintf("%s-%s-%s", tenantID, startDate.Format("200601"), endDate.Format("200601")),
		TenantID:          tenantID,
		ArchivePeriodStart: startDate,
		ArchivePeriodEnd:   endDate,
		LogCount:           int64(len(logs)),
		ArchiveLocation:    archiveURL,
		ArchiveHash:        hashHex,
		CreatedAt:          time.Now(),
	}
	
	// アーカイブメタデータを保存
	if err := s.archiveRepo.Create(ctx, archive); err != nil {
		return fmt.Errorf("failed to create archive metadata: %w", err)
	}
	
	// アーカイブしたログをDBから削除（またはarchived_atを設定）
	// 注意: 実際には削除せず、archived_atを設定してマークする方が安全
	if err := s.auditLogRepo.MarkAsArchived(ctx, tenantID, startDate, endDate, archiveURL); err != nil {
		return fmt.Errorf("failed to mark logs as archived: %w", err)
	}
	
	return nil
}

// VerifyArchiveIntegrity アーカイブファイルの整合性を検証
func (s *AuditArchiveService) VerifyArchiveIntegrity(ctx context.Context, archiveID string) (bool, error) {
	// アーカイブメタデータを取得
	// TODO: ArchiveRepositoryにGetByIDを追加する必要がある
	
	// Cloud Storageからアーカイブファイルをダウンロード
	// ハッシュ値を再計算
	// 保存されているハッシュ値と比較
	
	return false, fmt.Errorf("not implemented yet")
}

// GetArchiveLogs アーカイブからログを取得
func (s *AuditArchiveService) GetArchiveLogs(ctx context.Context, archiveID string) ([]*domain.AuditLog, error) {
	// アーカイブメタデータを取得
	// Cloud Storageからアーカイブファイルをダウンロード
	// JSONをデシリアライズ
	
	return nil, fmt.Errorf("not implemented yet")
}

