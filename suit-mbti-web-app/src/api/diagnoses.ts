import apiClient from './client';
import type {
  Diagnosis,
  CreateDiagnosisRequest,
  ListResponse,
} from '../types';

/**
 * 診断API
 */

/**
 * 診断を作成
 */
export const createDiagnosis = async (
  request: CreateDiagnosisRequest,
  tenantId: string
): Promise<Diagnosis> => {
  const response = await apiClient.post<Diagnosis>(
    `/api/diagnoses?tenant_id=${tenantId}`,
    request
  );
  return response.data;
};

/**
 * 診断を取得
 */
export const getDiagnosis = async (
  id: string,
  tenantId: string
): Promise<Diagnosis> => {
  const response = await apiClient.get<Diagnosis>(
    `/api/diagnoses/${id}?tenant_id=${tenantId}`
  );
  return response.data;
};

/**
 * 診断一覧を取得（テナント別）
 */
export const getDiagnosesByTenant = async (
  tenantId: string,
  options?: {
    limit?: number;
    offset?: number;
    archetype?: string;
    planType?: string;
  }
): Promise<ListResponse<Diagnosis>> => {
  const params = new URLSearchParams({
    tenant_id: tenantId,
  });

  if (options?.limit) {
    params.append('limit', options.limit.toString());
  }
  if (options?.offset) {
    params.append('offset', options.offset.toString());
  }
  if (options?.archetype) {
    params.append('archetype', options.archetype);
  }
  if (options?.planType) {
    params.append('plan_type', options.planType);
  }

  const response = await apiClient.get<ListResponse<Diagnosis>>(
    `/api/diagnoses?${params.toString()}`
  );
  return response.data;
};

/**
 * ユーザーの診断一覧を取得
 */
export const getDiagnosesByUser = async (
  userId: string,
  tenantId: string
): Promise<Diagnosis[]> => {
  const response = await apiClient.get<Diagnosis[]>(
    `/api/diagnoses?tenant_id=${tenantId}&user_id=${userId}`
  );
  return response.data;
};

/**
 * 診断を削除
 */
export const deleteDiagnosis = async (
  id: string,
  tenantId: string
): Promise<void> => {
  await apiClient.delete(`/api/diagnoses/${id}?tenant_id=${tenantId}`);
};

