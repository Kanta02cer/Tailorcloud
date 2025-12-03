/**
 * API設定
 * 環境変数からAPIベースURLを取得
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * APIクライアントのデフォルト設定
 */
export const API_CONFIG = {
  baseURL: API_BASE_URL,
  timeout: 30000, // 30秒
  headers: {
    'Content-Type': 'application/json',
  },
};

