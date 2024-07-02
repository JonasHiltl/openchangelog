package internal

type Pagination interface {
	PageSize() int
	Page() int
	StartIdx() int
	EndIdx() int
}

type pagination struct {
	pageSize int
	page     int
}

func NewPagination(pageSize int, page int) Pagination {
	return pagination{
		pageSize: pageSize,
		page:     page,
	}
}

func (p pagination) PageSize() int {
	return p.pageSize
}
func (p pagination) Page() int {
	return p.page
}

func (p pagination) StartIdx() int {
	return (p.page - 1) * p.pageSize
}

func (p pagination) EndIdx() int {
	return p.page*p.pageSize - 1
}
