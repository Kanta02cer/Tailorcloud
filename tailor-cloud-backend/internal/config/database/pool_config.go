package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PoolConfig データベース接続プール設定
type PoolConfig struct {
	MaxOpenConns    int           // 最大オープン接続数
	MaxIdleConns    int           // 最大アイドル接続数
	ConnMaxLifetime time.Duration // 接続の最大生存時間
	ConnMaxIdleTime time.Duration // アイドル接続の最大生存時間
}

// DefaultPoolConfig デフォルトの接続プール設定
// エンタープライズ要件: 100店舗×10工場×年間10万発注を想定
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,              // 最大25接続（Cloud SQL推奨値に基づく）
		MaxIdleConns:    10,              // 最大10アイドル接続
		ConnMaxLifetime: 5 * time.Minute, // 5分（Cloud SQL推奨値）
		ConnMaxIdleTime: 1 * time.Minute, // 1分
	}
}

// ConfigurePool データベース接続プールを設定
func ConfigurePool(db *sql.DB, config PoolConfig) error {
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}
	
	// 接続プールの状態を確認
	stats := db.Stats()
	fmt.Printf("Database connection pool configured: MaxOpenConns=%d, MaxIdleConns=%d, OpenConnections=%d, InUse=%d, Idle=%d\n",
		stats.MaxOpenConnections,
		config.MaxIdleConns,
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
	)
	
	return nil
}

// HighLoadPoolConfig 高負荷環境向けの接続プール設定
// 100店舗以上、同時接続数が多い場合
func HighLoadPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    50,              // 最大50接続
		MaxIdleConns:    20,              // 最大20アイドル接続
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

