// 注文およびコンプライアンス文書まわりの型定義

export type OrderStatus = 'Draft' | 'Confirmed' | 'Cancelled' | 'Completed' | string;

export interface Order {
  id: string;
  tenant_id: string;
  customer_id: string;
  fabric_id: string;
  status: OrderStatus;
  total_amount: number;
  tax_excluded_amount?: number;
  tax_amount?: number;
  tax_rate?: number;
  payment_due_date?: string;
  delivery_date: string;
  details?: {
    description?: string;
    measurement_data?: Record<string, unknown>;
    adjustments?: Record<string, unknown>;
  };
  created_at: string;
  updated_at: string;
  created_by?: string;
}

export interface CreateOrderRequest {
  customer_id: string;
  fabric_id: string; // バックエンドでは必須
  total_amount: number;
  delivery_date: string; // ISO 8601 (RFC3339形式: "2025-12-31T00:00:00Z")
  details: {
    description: string;
    measurement_data?: Record<string, unknown>;
    adjustments?: Record<string, unknown>;
  };
  // バックエンドで自動設定されるが、開発環境では明示的に指定可能
  tenant_id?: string;
  created_by?: string;
}

export interface GenerateDocumentResponse {
  order_id: string;
  doc_url: string;
  doc_hash: string;
  generated_at: string;
}


