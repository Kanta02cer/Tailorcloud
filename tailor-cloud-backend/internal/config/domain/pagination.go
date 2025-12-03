package domain

// Pagination ページネーション情報
type Pagination struct {
	Page     int `json:"page"`     // 現在のページ（1始まり）
	PageSize int `json:"page_size"` // 1ページあたりの件数
	Total    int `json:"total"`    // 全件数
}

// PaginatedResponse ページネーション付きレスポンス
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// NewPagination Paginationのコンストラクタ
func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20 // デフォルト: 20件
	}
	if pageSize > 100 {
		pageSize = 100 // 最大: 100件
	}
	
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// GetOffset オフセットを取得
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit リミットを取得
func (p Pagination) GetLimit() int {
	return p.PageSize
}

// GetTotalPages 総ページ数を取得
func (p Pagination) GetTotalPages() int {
	if p.Total == 0 {
		return 1
	}
	pages := p.Total / p.PageSize
	if p.Total%p.PageSize > 0 {
		pages++
	}
	return pages
}

// HasNext 次のページがあるか
func (p Pagination) HasNext() bool {
	return p.Page < p.GetTotalPages()
}

// HasPrev 前のページがあるか
func (p Pagination) HasPrev() bool {
	return p.Page > 1
}

