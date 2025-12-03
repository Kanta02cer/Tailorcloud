package service

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// StorageService Cloud Storageサービスのインターフェース
type StorageService interface {
	UploadPDF(ctx context.Context, bucketName string, objectPath string, pdfBytes []byte) (string, error)
	UploadJSON(ctx context.Context, bucketName string, objectPath string, jsonBytes []byte) (string, error)
	GetPublicURL(bucketName string, objectPath string) string
}

// GCSStorageService Google Cloud Storage実装
type GCSStorageService struct {
	client *storage.Client
}

// NewGCSStorageService GCSStorageServiceのコンストラクタ
func NewGCSStorageService(ctx context.Context, credentialsPath string) (*GCSStorageService, error) {
	var client *storage.Client
	var err error

	// 認証情報ファイルが指定されている場合は使用
	if credentialsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	} else {
		// 環境変数から認証情報を取得（GCP環境やサービスアカウント）
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &GCSStorageService{
		client: client,
	}, nil
}

// UploadPDF PDFファイルをCloud Storageにアップロード
func (s *GCSStorageService) UploadPDF(ctx context.Context, bucketName string, objectPath string, pdfBytes []byte) (string, error) {
	// バケットを取得
	bucket := s.client.Bucket(bucketName)

	// オブジェクトを作成
	obj := bucket.Object(objectPath)

	// アップロード処理（Contextタイムアウト設定）
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// ライターを取得
	writer := obj.NewWriter(ctx)
	writer.ContentType = "application/pdf"
	writer.CacheControl = "public, max-age=31536000" // 1年間キャッシュ

	// PDFバイトを書き込み
	if _, err := writer.Write(pdfBytes); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to write PDF bytes: %w", err)
	}

	// ライターを閉じる（アップロード完了）
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 公開URLを返す
	publicURL := s.GetPublicURL(bucketName, objectPath)
	return publicURL, nil
}

// UploadJSON JSONファイルをCloud Storageにアップロード（アーカイブ用）
func (s *GCSStorageService) UploadJSON(ctx context.Context, bucketName string, objectPath string, jsonBytes []byte) (string, error) {
	// バケットを取得
	bucket := s.client.Bucket(bucketName)

	// オブジェクトを作成
	obj := bucket.Object(objectPath)
	
	// WORM（Write Once Read Many）設定のため、オブジェクト保持ポリシーを設定
	// 注意: 実際のWORM設定はバケットレベルで行う必要がある
	writer := obj.NewWriter(ctx)
	writer.ContentType = "application/json"
	
	// JSONデータを書き込み
	if _, err := writer.Write(jsonBytes); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to write JSON to storage: %w", err)
	}
	
	// 書き込みをクローズ
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}
	
	// 公開URLを生成
	publicURL := s.GetPublicURL(bucketName, objectPath)
	
	return publicURL, nil
}

// GetPublicURL オブジェクトの公開URLを生成
func (s *GCSStorageService) GetPublicURL(bucketName string, objectPath string) string {
	// 公開URLの形式: https://storage.googleapis.com/{bucket}/{object}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
}

// Close クライアントを閉じる
func (s *GCSStorageService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

