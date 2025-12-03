package domain

// Permission 権限モデル
// リソースベースの細かい権限管理
type Permission struct {
	ID           string      `json:"id" db:"id"`
	TenantID     string      `json:"tenant_id" db:"tenant_id"`
	ResourceType ResourceType `json:"resource_type" db:"resource_type"`
	ResourceID   *string     `json:"resource_id" db:"resource_id"` // NULLの場合は全リソース
	Action       Action      `json:"action" db:"action"`
	Role         UserRole    `json:"role" db:"role"`
	Granted      bool        `json:"granted" db:"granted"`
}

// ResourceType リソースタイプ
type ResourceType string

const (
	ResourceTypeOrder             ResourceType = "ORDER"
	ResourceTypeCustomer          ResourceType = "CUSTOMER"
	ResourceTypeFabric            ResourceType = "FABRIC"
	ResourceTypeInvoice           ResourceType = "INVOICE"
	ResourceTypeComplianceDocument ResourceType = "COMPLIANCE_DOCUMENT"
	ResourceTypeFabricRoll        ResourceType = "FABRIC_ROLL"
	ResourceTypeAll               ResourceType = "ALL"
)

// Action 操作タイプ
type Action string

const (
	ActionCreate    Action = "CREATE"
	ActionRead      Action = "READ"
	ActionUpdate    Action = "UPDATE"
	ActionDelete    Action = "DELETE"
	ActionGenerate  Action = "GENERATE"  // PDF生成など
	ActionApprove   Action = "APPROVE"   // 承認
	ActionView      Action = "VIEW"      // 閲覧のみ
	ActionAll       Action = "ALL"       // 全操作
)

// PermissionCheck 権限チェック結果
type PermissionCheck struct {
	Allowed  bool
	Reason   string
	Permission *Permission
}

// CheckPermission 権限をチェック
// resourceType: リソースタイプ
// action: 操作
// role: ユーザーロール
// resourceID: 特定リソースID（オプション）
func (p *Permission) CheckPermission(resourceType ResourceType, action Action, role UserRole, resourceID *string) bool {
	// ロールが一致しない場合は許可しない
	if p.Role != role {
		return false
	}
	
	// 全リソースに対する権限の場合
	if p.ResourceType == ResourceTypeAll {
		return p.Granted
	}
	
	// リソースタイプが一致するか
	if p.ResourceType != resourceType {
		return false
	}
	
	// 特定リソースに対する権限の場合、IDが一致する必要がある
	if p.ResourceID != nil {
		if resourceID == nil || *p.ResourceID != *resourceID {
			return false
		}
	}
	
	// 全操作に対する権限の場合
	if p.Action == ActionAll {
		return p.Granted
	}
	
	// 操作が一致するか
	if p.Action != action {
		return false
	}
	
	return p.Granted
}

