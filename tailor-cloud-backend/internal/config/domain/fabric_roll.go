package domain

import "time"

// FabricRollStatus 反物（Roll）の状態
type FabricRollStatus string

const (
	FabricRollStatusAvailable FabricRollStatus = "AVAILABLE" // 利用可能
	FabricRollStatusAllocated FabricRollStatus = "ALLOCATED" // 引当済み
	FabricRollStatusConsumed  FabricRollStatus = "CONSUMED"  // 消費済み
	FabricRollStatusDamaged   FabricRollStatus = "DAMAGED"   // 破損
)

// FabricRoll 反物（Roll）モデル
// 物理的な1本の反物を表す。単なる「総量」ではなく「物理的な巻き」単位で管理
type FabricRoll struct {
	ID            string           `json:"id"`
	TenantID      string           `json:"tenant_id"`
	FabricID      string           `json:"fabric_id"`
	RollNumber    string           `json:"roll_number"`    // ロール番号（例: "VBC-2025-001"）
	InitialLength float64          `json:"initial_length"` // 初期長さ（メートル）
	CurrentLength float64          `json:"current_length"` // 現在の残り長さ（メートル）
	Width         *float64         `json:"width,omitempty"` // 幅（センチメートル、オプション）
	SupplierLotNo *string          `json:"supplier_lot_no,omitempty"` // 仕入先ロット番号
	ReceivedAt    *time.Time       `json:"received_at,omitempty"`     // 入荷日
	Location      *string          `json:"location,omitempty"`        // 保管場所
	Status        FabricRollStatus `json:"status"`                    // 状態
	Notes         *string          `json:"notes,omitempty"`           // 備考
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// CanAllocate 引当可能かどうかをチェック
func (r *FabricRoll) CanAllocate(requiredLength float64) bool {
	return r.Status == FabricRollStatusAvailable && r.CurrentLength >= requiredLength
}

// Allocate 反物を引当（残り長さを減算）
func (r *FabricRoll) Allocate(allocatedLength float64) error {
	if !r.CanAllocate(allocatedLength) {
		return ErrCannotAllocate
	}
	
	if r.CurrentLength < allocatedLength {
		return ErrInsufficientLength
	}
	
	r.CurrentLength -= allocatedLength
	if r.CurrentLength == 0 {
		r.Status = FabricRollStatusConsumed
	} else {
		r.Status = FabricRollStatusAllocated
	}
	
	r.UpdatedAt = time.Now()
	return nil
}

// Release 引当を解除（残り長さを戻す）
func (r *FabricRoll) Release(releasedLength float64) {
	r.CurrentLength += releasedLength
	
	if r.Status == FabricRollStatusAllocated && r.CurrentLength > 0 {
		// 他の引当が残っている可能性があるため、ステータスはALLOCATEDのまま
	} else if r.Status == FabricRollStatusConsumed {
		// 消費済みから戻す場合はAVAILABLEに
		r.Status = FabricRollStatusAvailable
	}
	
	r.UpdatedAt = time.Now()
}

// FabricAllocationStatus 引当状態
type FabricAllocationStatus string

const (
	FabricAllocationStatusReserved  FabricAllocationStatus = "RESERVED"  // 予約済み
	FabricAllocationStatusConfirmed FabricAllocationStatus = "CONFIRMED" // 確定
	FabricAllocationStatusCut       FabricAllocationStatus = "CUT"       // 裁断済み
	FabricAllocationStatusCancelled FabricAllocationStatus = "CANCELLED" // キャンセル
)

// FabricAllocation 反物引当モデル
// 発注時にどの反物（Roll）のどの部分を使用するかを記録
type FabricAllocation struct {
	ID               string                 `json:"id"`
	TenantID         string                 `json:"tenant_id"`
	OrderID          string                 `json:"order_id"`
	OrderItemID      *string                `json:"order_item_id,omitempty"` // 将来の拡張用
	FabricRollID     string                 `json:"fabric_roll_id"`
	AllocatedLength  float64                `json:"allocated_length"`  // 引当数量（メートル）
	ActualUsedLength *float64               `json:"actual_used_length,omitempty"` // 実際に使用した数量
	RemnantLength    *float64               `json:"remnant_length,omitempty"`     // 端尺（キレ）の長さ
	Status           FabricAllocationStatus `json:"status"`                       // 引当状態
	AllocatedAt      time.Time              `json:"allocated_at"`                 // 引当日時
	ConfirmedAt      *time.Time             `json:"confirmed_at,omitempty"`       // 確定日時
	CutAt            *time.Time             `json:"cut_at,omitempty"`             // 裁断日時
	Notes            *string                `json:"notes,omitempty"`              // 備考
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// Confirm 引当を確定
func (a *FabricAllocation) Confirm() {
	now := time.Now()
	a.Status = FabricAllocationStatusConfirmed
	a.ConfirmedAt = &now
	a.UpdatedAt = now
}

// MarkAsCut 裁断済みとしてマーク
func (a *FabricAllocation) MarkAsCut(actualUsedLength, remnantLength float64) {
	now := time.Now()
	a.Status = FabricAllocationStatusCut
	a.ActualUsedLength = &actualUsedLength
	a.RemnantLength = &remnantLength
	a.CutAt = &now
	a.UpdatedAt = now
}

// Cancel 引当をキャンセル
func (a *FabricAllocation) Cancel() {
	now := time.Now()
	a.Status = FabricAllocationStatusCancelled
	a.UpdatedAt = now
}

// エラー定義
var (
	ErrCannotAllocate    = NewDomainError("fabric roll cannot be allocated", "FABRIC_ROLL_CANNOT_ALLOCATE")
	ErrInsufficientLength = NewDomainError("insufficient fabric roll length", "FABRIC_ROLL_INSUFFICIENT_LENGTH")
)

// DomainError ドメインエラー
type DomainError struct {
	Message string
	Code    string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(message, code string) *DomainError {
	return &DomainError{
		Message: message,
		Code:    code,
	}
}

