package service

import (
	"context"
	"fmt"
	"time"

	"tailor-cloud/backend/internal/repository"
)

// AnalyticsSummaryRequest リクエストペイロード
type AnalyticsSummaryRequest struct {
	TenantID  string
	RangeDays int
}

// TagCount タグ別集計
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// AnalyticsSummary KPIレスポンス
type AnalyticsSummary struct {
	TenantID          string         `json:"tenant_id"`
	RangeDays         int            `json:"range_days"`
	GeneratedAt       time.Time      `json:"generated_at"`
	TotalOrders       int            `json:"total_orders"`
	TotalRevenue      int64          `json:"total_revenue"`
	AverageOrderValue float64        `json:"average_order_value"`
	ActiveCustomers   int            `json:"active_customers"`
	StatusBreakdown   map[string]int `json:"status_breakdown"`
	TopTags           []TagCount     `json:"top_tags"`
}

// AnalyticsService ダッシュボード用KPIサービス
type AnalyticsService struct {
	orderRepo    repository.OrderRepository
	customerRepo repository.CustomerRepository
}

// NewAnalyticsService コンストラクタ
func NewAnalyticsService(orderRepo repository.OrderRepository, customerRepo repository.CustomerRepository) *AnalyticsService {
	return &AnalyticsService{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
	}
}

// GetSummary KPIサマリーを計算
func (s *AnalyticsService) GetSummary(ctx context.Context, req *AnalyticsSummaryRequest) (*AnalyticsSummary, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if req.RangeDays <= 0 {
		req.RangeDays = 30
	}

	orders, err := s.orderRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	customers, err := s.customerRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customers: %w", err)
	}

	fromDate := time.Now().AddDate(0, 0, -req.RangeDays)
	var (
		totalOrders  int
		totalRevenue int64
	)

	for _, order := range orders {
		if order.CreatedAt.Before(fromDate) {
			continue
		}
		totalOrders++
		totalRevenue += order.TotalAmount
	}

	var average float64
	if totalOrders > 0 && totalRevenue > 0 {
		average = float64(totalRevenue) / float64(totalOrders)
	}

	statusBreakdown := make(map[string]int)
	tagCounter := make(map[string]int)
	activeCustomers := 0

	for _, customer := range customers {
		status := customer.CustomerStatus
		if status == "" {
			status = "lead"
		}
		statusBreakdown[status]++

		if status == "active" || status == "vip" {
			activeCustomers++
		}

		for _, tag := range customer.Tags {
			tagCounter[tag]++
		}
	}

	topTags := buildTopTags(tagCounter, 5)

	return &AnalyticsSummary{
		TenantID:          req.TenantID,
		RangeDays:         req.RangeDays,
		GeneratedAt:       time.Now(),
		TotalOrders:       totalOrders,
		TotalRevenue:      totalRevenue,
		AverageOrderValue: average,
		ActiveCustomers:   activeCustomers,
		StatusBreakdown:   statusBreakdown,
		TopTags:           topTags,
	}, nil
}

func buildTopTags(counter map[string]int, limit int) []TagCount {
	results := make([]TagCount, 0, len(counter))
	for tag, count := range counter {
		results = append(results, TagCount{Tag: tag, Count: count})
	}
	if len(results) <= 1 {
		return results
	}
	// simple selection sort for small dataset
	for i := 0; i < len(results); i++ {
		maxIdx := i
		for j := i + 1; j < len(results); j++ {
			if results[j].Count > results[maxIdx].Count {
				maxIdx = j
			}
		}
		results[i], results[maxIdx] = results[maxIdx], results[i]
	}
	if len(results) > limit {
		return results[:limit]
	}
	return results
}
