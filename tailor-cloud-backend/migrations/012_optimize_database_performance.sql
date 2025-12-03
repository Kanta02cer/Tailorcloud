-- ============================================================================
-- TailorCloud Enterprise: データベースパフォーマンス最適化
-- ============================================================================
-- 目的: 100店舗×10工場×年間10万発注の処理が可能に
-- インデックス最適化、複合インデックス追加、クエリパフォーマンス向上
-- ============================================================================

-- 1. 注文テーブルの最適化
-- よく使われるクエリパターンに基づく複合インデックス

-- テナント×ステータス×作成日の複合インデックス（ダッシュボード表示用）
CREATE INDEX IF NOT EXISTS idx_orders_tenant_status_created 
ON orders(tenant_id, status, created_at DESC);

-- テナント×顧客ID×作成日の複合インデックス（顧客別注文一覧用）
CREATE INDEX IF NOT EXISTS idx_orders_tenant_customer_created 
ON orders(tenant_id, customer_id, created_at DESC);

-- テナント×配送日の複合インデックス（配送スケジュール確認用）
CREATE INDEX IF NOT EXISTS idx_orders_tenant_delivery_date 
ON orders(tenant_id, delivery_date) WHERE delivery_date >= CURRENT_DATE;

-- 2. 顧客テーブルの最適化
-- テナント×作成日の複合インデックス（ページネーション用）
CREATE INDEX IF NOT EXISTS idx_customers_tenant_created 
ON customers(tenant_id, created_at DESC);

-- 3. 生地テーブルの最適化
-- 検索パフォーマンス向上（ILIKE検索の最適化）
-- 注意: 完全一致検索には効果的だが、前方一致検索には別途トリグラムインデックスが必要
CREATE INDEX IF NOT EXISTS idx_fabrics_name_lower 
ON fabrics(LOWER(name));

-- 4. 反物管理テーブルの最適化
-- テナント×生地ID×ステータス×現在長の複合インデックス（在庫検索用）
CREATE INDEX IF NOT EXISTS idx_fabric_rolls_tenant_fabric_status_length 
ON fabric_rolls(tenant_id, fabric_id, status, current_length DESC) 
WHERE status IN ('AVAILABLE', 'ALLOCATED');

-- 5. 在庫引当テーブルの最適化
-- テナント×注文ID×ステータスの複合インデックス（注文別引当状況確認用）
CREATE INDEX IF NOT EXISTS idx_fabric_allocations_tenant_order_status 
ON fabric_allocations(tenant_id, order_id, allocation_status);

-- 6. アンバサダー成果報酬テーブルの最適化
-- テナント×アンバサダーID×ステータス×作成日の複合インデックス（成果報酬一覧用）
CREATE INDEX IF NOT EXISTS idx_commissions_tenant_ambassador_status_created 
ON commissions(tenant_id, ambassador_id, status, created_at DESC);

-- 7. 監査ログテーブルの最適化（既存のインデックスを補完）
-- テナント×リソースタイプ×アクション×作成日の複合インデックス（監査検索用）
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_resource_action_created 
ON audit_logs(tenant_id, resource_type, action, created_at DESC);

-- 8. 部分インデックス（条件付きインデックス）による最適化
-- アクティブな注文のみをインデックス（ステータスがDraft/Confirmed等の注文）
CREATE INDEX IF NOT EXISTS idx_orders_active_status 
ON orders(tenant_id, created_at DESC) 
WHERE status IN ('Draft', 'Confirmed', 'Material_Secured', 'Cutting', 'Sewing', 'Inspection');

-- 9. カバリングインデックス（必要なカラムのみを含むインデックス）の検討
-- 注：PostgreSQL 11以降でINCLUDE句が使用可能な場合、検索対象カラムのみを含むインデックスを検討

-- 10. 分析クエリ用の統計情報更新
-- PostgreSQLの自動統計情報収集を促進するための設定（コメント）
COMMENT ON INDEX idx_orders_tenant_status_created IS 'ダッシュボード表示用（テナント×ステータス×作成日）';
COMMENT ON INDEX idx_orders_tenant_customer_created IS '顧客別注文一覧用（テナント×顧客×作成日）';
COMMENT ON INDEX idx_customers_tenant_created IS '顧客一覧ページネーション用（テナント×作成日）';

-- 11. パフォーマンス分析用のビュー（オプション）
-- 注文統計ビュー（分析用）
CREATE OR REPLACE VIEW v_order_stats_by_tenant AS
SELECT 
    tenant_id,
    status,
    COUNT(*) as order_count,
    SUM(total_amount) as total_amount_sum,
    AVG(total_amount) as avg_amount,
    MIN(created_at) as earliest_order,
    MAX(created_at) as latest_order
FROM orders
GROUP BY tenant_id, status;

COMMENT ON VIEW v_order_stats_by_tenant IS 'テナント別注文統計ビュー（分析用）';

-- 12. インデックス使用状況の監視（コメント）
-- PostgreSQLのpg_stat_user_indexesビューを使用してインデックスの使用状況を監視
-- 未使用のインデックスは削除を検討（定期的なメンテナンスが必要）

