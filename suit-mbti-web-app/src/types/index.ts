// TypeScript型定義

// アーキタイプ（RATタイプ）
export type Archetype = 'Classic' | 'Modern' | 'Elegant' | 'Sporty' | 'Casual';

// プランタイプ
export type PlanType = 'Best Value' | 'Authentic';

// 診断結果
export interface Diagnosis {
  id: string;
  user_id: string;
  tenant_id: string;
  archetype: Archetype;
  plan_type?: PlanType;
  diagnosis_result: DiagnosisResult;
  created_at: string;
  updated_at: string;
}

// 診断結果詳細
export interface DiagnosisResult {
  scores: {
    classic?: number;
    modern?: number;
    elegant?: number;
    sporty?: number;
    casual?: number;
  };
  recommendations?: string[];
  notes?: string;
}

// 診断作成リクエスト
export interface CreateDiagnosisRequest {
  user_id: string;
  archetype: Archetype;
  plan_type?: PlanType;
  diagnosis_result: DiagnosisResult;
}

// 予約ステータス
export type AppointmentStatus = 'Pending' | 'Confirmed' | 'Cancelled' | 'Completed' | 'NoShow';

// デポジットステータス
export type DepositStatus = 'pending' | 'succeeded' | 'failed' | 'refunded' | '';

// 予約
export interface Appointment {
  id: string;
  user_id: string;
  tenant_id: string;
  fitter_id?: string;
  appointment_datetime: string;
  duration_minutes: number;
  status: AppointmentStatus;
  deposit_amount?: number | null;
  deposit_payment_intent_id?: string;
  deposit_status?: DepositStatus;
  notes?: string;
  cancelled_at?: string | null;
  cancelled_reason?: string;
  created_at: string;
  updated_at: string;
}

// 予約作成リクエスト
export interface CreateAppointmentRequest {
  user_id: string;
  fitter_id?: string;
  appointment_datetime: string; // ISO 8601 format (RFC3339)
  duration_minutes: number;
  notes?: string;
}

// 予約更新リクエスト
export interface UpdateAppointmentRequest {
  fitter_id?: string;
  appointment_datetime?: string;
  duration_minutes?: number;
  status?: AppointmentStatus;
  notes?: string;
}

// APIレスポンス（リスト取得）
export interface ListResponse<T> {
  data: T[];
  total: number;
}

