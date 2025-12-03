package service

import (
	"testing"

	"tailor-cloud/backend/internal/config/domain"
)

// TestInventoryAllocationService_AllocateInventory 在庫引当サービスのテスト
// 注意: 実際のPostgreSQL接続が必要（統合テスト）
func TestInventoryAllocationService_AllocateInventory(t *testing.T) {
	// このテストは統合テストとして実行
	// 実際のPostgreSQL接続が必要
	// TODO: テスト用DBのセットアップが必要
	t.Skip("Skipping integration test - requires PostgreSQL connection")
}

// TestFabricRoll_CanAllocate 反物の引当可能性チェックのテスト
func TestFabricRoll_CanAllocate(t *testing.T) {
	roll := &domain.FabricRoll{
		ID:            "test-roll-1",
		TenantID:      "test-tenant",
		FabricID:      "test-fabric",
		RollNumber:    "TEST-001",
		InitialLength: 50.0,
		CurrentLength: 10.0,
		Status:        domain.FabricRollStatusAvailable,
	}

	// テストケース1: 引当可能（残り長さ >= 必要長さ）
	if !roll.CanAllocate(5.0) {
		t.Error("Expected CanAllocate(5.0) to return true, got false")
	}

	// テストケース2: 引当不可（残り長さ < 必要長さ）
	if roll.CanAllocate(15.0) {
		t.Error("Expected CanAllocate(15.0) to return false, got true")
	}

	// テストケース3: 引当不可（ステータスがAVAILABLEではない）
	roll.Status = domain.FabricRollStatusAllocated
	if roll.CanAllocate(5.0) {
		t.Error("Expected CanAllocate(5.0) to return false when status is ALLOCATED, got true")
	}
}

// TestFabricRoll_Allocate 反物の引当処理のテスト
func TestFabricRoll_Allocate(t *testing.T) {
	roll := &domain.FabricRoll{
		ID:            "test-roll-1",
		TenantID:      "test-tenant",
		FabricID:      "test-fabric",
		RollNumber:    "TEST-001",
		InitialLength: 50.0,
		CurrentLength: 10.0,
		Status:        domain.FabricRollStatusAvailable,
	}

	// テストケース1: 正常な引当
	err := roll.Allocate(5.0)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if roll.CurrentLength != 5.0 {
		t.Errorf("Expected CurrentLength to be 5.0, got: %.2f", roll.CurrentLength)
	}
	if roll.Status != domain.FabricRollStatusAllocated {
		t.Errorf("Expected Status to be ALLOCATED, got: %s", roll.Status)
	}

	// テストケース2: 残り長さが0になる場合
	roll.Status = domain.FabricRollStatusAvailable // 再度AVAILABLEに戻す
	err = roll.Allocate(5.0)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if roll.CurrentLength != 0.0 {
		t.Errorf("Expected CurrentLength to be 0.0, got: %.2f", roll.CurrentLength)
	}
	if roll.Status != domain.FabricRollStatusConsumed {
		t.Errorf("Expected Status to be CONSUMED, got: %s", roll.Status)
	}

	// テストケース3: 引当不可（残り長さ不足）
	roll.CurrentLength = 2.0
	roll.Status = domain.FabricRollStatusAvailable
	err = roll.Allocate(5.0)
	if err == nil {
		t.Error("Expected error when allocating more than available, got nil")
	}
}

// TestFabricRoll_Release 反物の引当解除のテスト
func TestFabricRoll_Release(t *testing.T) {
	roll := &domain.FabricRoll{
		ID:            "test-roll-1",
		TenantID:      "test-tenant",
		FabricID:      "test-fabric",
		RollNumber:    "TEST-001",
		InitialLength: 50.0,
		CurrentLength: 5.0,
		Status:        domain.FabricRollStatusAllocated,
	}

	// 引当解除
	roll.Release(5.0)
	if roll.CurrentLength != 10.0 {
		t.Errorf("Expected CurrentLength to be 10.0, got: %.2f", roll.CurrentLength)
	}

	// 消費済みから戻す場合
	roll.CurrentLength = 0.0
	roll.Status = domain.FabricRollStatusConsumed
	roll.Release(5.0)
	if roll.Status != domain.FabricRollStatusAvailable {
		t.Errorf("Expected Status to be AVAILABLE after releasing from CONSUMED, got: %s", roll.Status)
	}
}

// TestFabricAllocation_Confirm 引当確定のテスト
func TestFabricAllocation_Confirm(t *testing.T) {
	allocation := &domain.FabricAllocation{
		ID:              "test-allocation-1",
		TenantID:        "test-tenant",
		OrderID:         "test-order",
		FabricRollID:    "test-roll",
		AllocatedLength: 3.2,
		Status:          domain.FabricAllocationStatusReserved,
	}

	allocation.Confirm()
	if allocation.Status != domain.FabricAllocationStatusConfirmed {
		t.Errorf("Expected Status to be CONFIRMED, got: %s", allocation.Status)
	}
	if allocation.ConfirmedAt == nil {
		t.Error("Expected ConfirmedAt to be set, got nil")
	}
}

// TestFabricAllocation_MarkAsCut 裁断済みマークのテスト
func TestFabricAllocation_MarkAsCut(t *testing.T) {
	allocation := &domain.FabricAllocation{
		ID:              "test-allocation-1",
		TenantID:        "test-tenant",
		OrderID:         "test-order",
		FabricRollID:    "test-roll",
		AllocatedLength: 3.2,
		Status:          domain.FabricAllocationStatusConfirmed,
	}

	actualUsed := 3.0
	remnant := 0.2
	allocation.MarkAsCut(actualUsed, remnant)

	if allocation.Status != domain.FabricAllocationStatusCut {
		t.Errorf("Expected Status to be CUT, got: %s", allocation.Status)
	}
	if allocation.ActualUsedLength == nil || *allocation.ActualUsedLength != actualUsed {
		t.Errorf("Expected ActualUsedLength to be %.2f, got: %v", actualUsed, allocation.ActualUsedLength)
	}
	if allocation.RemnantLength == nil || *allocation.RemnantLength != remnant {
		t.Errorf("Expected RemnantLength to be %.2f, got: %v", remnant, allocation.RemnantLength)
	}
	if allocation.CutAt == nil {
		t.Error("Expected CutAt to be set, got nil")
	}
}

