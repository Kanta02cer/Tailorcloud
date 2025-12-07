package middleware

import (
	"net/http"
)

// CORSMiddleware CORSミドルウェア
// 開発環境ではすべてのオリジンを許可
type CORSMiddleware struct {
	allowedOrigins []string
}

// NewCORSMiddleware CORSMiddlewareのコンストラクタ
func NewCORSMiddleware(allowedOrigins []string) *CORSMiddleware {
	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Handle CORSヘッダーを設定するミドルウェア
func (m *CORSMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// 開発環境ではすべてのオリジンを許可
		// 本番環境ではallowedOriginsをチェック
		allowedOrigin := "*"
		if len(m.allowedOrigins) > 0 {
			// 許可されたオリジンかチェック
			for _, allowed := range m.allowedOrigins {
				if origin == allowed || allowed == "*" {
					allowedOrigin = origin
					break
				}
			}
		}

		// CORSヘッダーを設定
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// OPTIONSリクエスト（プリフライト）の処理
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 次のハンドラーへ
		next.ServeHTTP(w, r)
	}
}

// CORS CORSヘッダーを設定する簡易ミドルウェア（開発環境用）
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// CORSヘッダーを設定
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// OPTIONSリクエスト（プリフライト）の処理
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 次のハンドラーへ
		next.ServeHTTP(w, r)
	}
}

