// 顧客・注文まわりの型定義

export interface Customer {
  id: string;
  tenant_id: string;
  name: string;
  email?: string;
  phone?: string;
  address?: string;
  created_at: string;
  updated_at: string;
}

export interface CustomerListResponse {
  customers: Customer[];
  total: number;
}

// 注文ステータス（既存仕様を踏まえつつ柔軟に）
export type OrderStatus = 'Draft' | 'Confirmed' | 'Cancelled' | 'Completed' | string;

// 顧客詳細画面で利用する注文サマリー
export interface OrderSummary {
  id: string;
  tenant_id: string;
  customer_id: string;
  status: OrderStatus;
  total_amount: number;
  delivery_date?: string;
  created_at: string;
}


