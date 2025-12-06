package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	dbpool "tailor-cloud/backend/internal/config/database"
)

// DatabaseConfig データベース設定
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDatabaseConfig 環境変数からデータベース設定を読み込む
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "tailorcloud"),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		DBName:   getEnv("POSTGRES_DB", "tailorcloud"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"), // Cloud SQLの場合は "disable" または "require"
	}
}

// getEnv 環境変数を取得（デフォルト値あり）
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// BuildConnectionString 接続文字列を構築
// DATABASE_URL環境変数が設定されている場合はそれを優先的に使用
func (c *DatabaseConfig) BuildConnectionString() string {
	// DATABASE_URL環境変数が設定されている場合はそれを使用（Render, Heroku等のクラウドプラットフォーム対応）
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}
	
	// DATABASE_URLが設定されていない場合は個別の環境変数から構築
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// ConnectPostgreSQL PostgreSQLに接続
func ConnectPostgreSQL(config *DatabaseConfig) (*sql.DB, error) {
	connStr := config.BuildConnectionString()
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// コネクションプール設定（デフォルト設定）
	if err := dbpool.ConfigurePool(db, dbpool.DefaultPoolConfig()); err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}
	
	return db, nil
}

// ConfigurePool 接続プールを設定（外部から呼び出し可能）
func ConfigurePool(db *sql.DB, config dbpool.PoolConfig) error {
	return dbpool.ConfigurePool(db, config)
}

// DefaultPoolConfig デフォルトの接続プール設定を取得
func DefaultPoolConfig() dbpool.PoolConfig {
	return dbpool.DefaultPoolConfig()
}

// HighLoadPoolConfig 高負荷環境向けの接続プール設定を取得
func HighLoadPoolConfig() dbpool.PoolConfig {
	return dbpool.HighLoadPoolConfig()
}

