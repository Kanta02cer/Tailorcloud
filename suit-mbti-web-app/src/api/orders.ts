import apiClient from './client';
import type { CreateOrderRequest, GenerateDocumentResponse, Order } from '../types/order';

const DEFAULT_TENANT_ID = import.meta.env.VITE_TENANT_ID || 'tenant_test_suit_mbti';

export const createOrder = async (request: CreateOrderRequest): Promise<Order> => {
  // バックエンドの仕様に合わせて、tenant_idを追加（開発環境用）
  const payload: CreateOrderRequest & { tenant_id?: string } = {
    ...request,
    tenant_id: request.tenant_id || DEFAULT_TENANT_ID,
  };
  const response = await apiClient.post<Order>('/api/orders', payload);
  return response.data;
};

export const generateOrderDocument = async (orderId: string): Promise<GenerateDocumentResponse> => {
  const response = await apiClient.post<GenerateDocumentResponse>(
    `/api/orders/${orderId}/generate-document`
  );
  return response.data;
};


