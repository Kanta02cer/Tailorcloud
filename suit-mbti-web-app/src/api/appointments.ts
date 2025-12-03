import apiClient from './client';
import type {
  Appointment,
  CreateAppointmentRequest,
  UpdateAppointmentRequest,
  ListResponse,
} from '../types';

/**
 * 予約API
 */

/**
 * 予約を作成
 */
export const createAppointment = async (
  request: CreateAppointmentRequest,
  tenantId: string
): Promise<Appointment> => {
  const response = await apiClient.post<Appointment>(
    `/api/appointments?tenant_id=${tenantId}`,
    request
  );
  return response.data;
};

/**
 * 予約を取得
 */
export const getAppointment = async (
  id: string,
  tenantId: string
): Promise<Appointment> => {
  const response = await apiClient.get<Appointment>(
    `/api/appointments/${id}?tenant_id=${tenantId}`
  );
  return response.data;
};

/**
 * 予約一覧を取得
 */
export const listAppointments = async (
  tenantId: string,
  options?: {
    userId?: string;
    fitterId?: string;
    startDate?: string; // ISO 8601 format
    endDate?: string; // ISO 8601 format
  }
): Promise<ListResponse<Appointment>> => {
  const params = new URLSearchParams({
    tenant_id: tenantId,
  });

  if (options?.userId) {
    params.append('user_id', options.userId);
  }
  if (options?.fitterId) {
    params.append('fitter_id', options.fitterId);
  }
  if (options?.startDate) {
    params.append('start_date', options.startDate);
  }
  if (options?.endDate) {
    params.append('end_date', options.endDate);
  }

  const response = await apiClient.get<ListResponse<Appointment>>(
    `/api/appointments?${params.toString()}`
  );
  return response.data;
};

/**
 * 予約を更新
 */
export const updateAppointment = async (
  id: string,
  request: UpdateAppointmentRequest,
  tenantId: string
): Promise<Appointment> => {
  const response = await apiClient.patch<Appointment>(
    `/api/appointments/${id}?tenant_id=${tenantId}`,
    request
  );
  return response.data;
};

/**
 * 予約をキャンセル
 */
export const cancelAppointment = async (
  id: string,
  reason: string,
  tenantId: string
): Promise<void> => {
  await apiClient.delete(
    `/api/appointments/${id}?tenant_id=${tenantId}&reason=${encodeURIComponent(reason)}`
  );
};

/**
 * 空き状況を確認
 */
export const checkAvailability = async (
  fitterId: string,
  date: string, // ISO 8601 format
  tenantId: string
): Promise<{ available: boolean; slots: string[] }> => {
  const response = await apiClient.get<{ available: boolean; slots: string[] }>(
    `/api/appointments/availability?tenant_id=${tenantId}&fitter_id=${fitterId}&date=${date}`
  );
  return response.data;
};

