import axios, { AxiosInstance, AxiosError } from 'axios';

// 環境変数からAPIベースURLを取得
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Axiosインスタンスを作成
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエストインターセプター
apiClient.interceptors.request.use(
  (config) => {
    // 認証トークンがあれば追加（今後実装）
    // const token = getAuthToken();
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // エラーハンドリング（詳細なメッセージを返す）
    if (error.response) {
      // サーバーからのエラーレスポンス
      const status = error.response.status;
      let message: string = 'エラーが発生しました';
      
      // バックエンドからのエラーメッセージを抽出
      if (typeof error.response.data === 'string') {
        message = error.response.data;
      } else if (error.response.data && typeof error.response.data === 'object') {
        const data = error.response.data as { error?: string; message?: string };
        message = data.error || data.message || message;
      }
      
      // エラーオブジェクトに詳細情報を追加
      const enhancedError = new Error(message);
      (enhancedError as any).status = status;
      (enhancedError as any).response = error.response;
      
      switch (status) {
        case 400:
          console.error('Bad Request:', message);
          break;
        case 401:
          console.error('Unauthorized: Please check your authentication');
          break;
        case 403:
          console.error('Forbidden: You do not have permission to access this resource');
          break;
        case 404:
          console.error('Not Found: The requested resource was not found');
          break;
        case 409:
          console.error('Conflict:', message);
          break;
        case 500:
          console.error('Internal Server Error: Please try again later');
          break;
        default:
          console.error(`Error ${status}:`, message);
      }
      
      return Promise.reject(enhancedError);
    } else if (error.request) {
      // リクエストは送信されたが、レスポンスが返ってこなかった
      const networkError = new Error(
        'ネットワークエラー: インターネット接続を確認してください'
      );
      (networkError as any).isNetworkError = true;
      console.error('Network Error: Please check your internet connection');
      return Promise.reject(networkError);
    } else {
      // リクエスト設定時にエラーが発生
      console.error('Error:', error.message);
      return Promise.reject(error);
    }
  }
);

export default apiClient;

