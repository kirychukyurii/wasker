package dto

type Pagination struct {
	Total    uint64 `json:"total"`
	Current  uint64 `json:"current"`
	PageSize uint64 `json:"page_size"`
}

type PaginationParam struct {
	Current  uint64 `json:"current" query:"current"`
	PageSize uint64 `json:"page_size" query:"page_size" validate:"max=128"`
}

func (a *PaginationParam) GetCurrent() uint64 {
	return a.Current
}

func (a *PaginationParam) GetPageSize() uint64 {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 15
	}

	return pageSize
}
