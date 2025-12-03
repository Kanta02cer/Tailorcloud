import apiClient from './client';
import type { Customer, CustomerListResponse, OrderSummary } from '../types/customer';

const DEFAULT_TENANT_ID = import.meta.env.VITE_TENANT_ID || 'tenant_test_suit_mbti';

export interface ListCustomersOptions {
  search?: string;
  tenantId?: string;
}

export const listCustomers = async (
  options: ListCustomersOptions = {}
): Promise<CustomerListResponse> => {
  const params = new URLSearchParams();
  const tenantId = options.tenantId ?? DEFAULT_TENANT_ID;

  params.append('tenant_id', tenantId);
  if (options.search) {
    params.append('search', options.search);
  }

  const response = await apiClient.get<CustomerListResponse>(`/api/customers?${params.toString()}`);
  return response.data;
};

export const getCustomer = async (
  id: string,
  tenantId: string = DEFAULT_TENANT_ID
): Promise<Customer> => {
  const response = await apiClient.get<Customer>(`/api/customers/${id}?tenant_id=${tenantId}`);
  return response.data;
};

export const getCustomerOrders = async (
  id: string,
  tenantId: string = DEFAULT_TENANT_ID
): Promise<OrderSummary[]> => {
  const response = await apiClient.get<OrderSummary[]>(
    `/api/customers/${id}/orders?tenant_id=${tenantId}`
  );
  return response.data;
};


