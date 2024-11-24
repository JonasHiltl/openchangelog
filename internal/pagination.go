package internal

type Pagination interface {
	PageSize() int
	Page() int
	StartIdx() int
	EndIdx() int
	// Returns true if the pagination is defined, else false and pagination should be ignored
	IsDefined() bool
}

type pagination struct {
	pageSize  int
	page      int
	isDefined bool
}

func NewPagination(pageSize int, page int) Pagination {
	return pagination{
		pageSize:  pageSize,
		page:      page,
		isDefined: true,
	}
}

func NoPagination() Pagination {
	return pagination{
		isDefined: false,
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

func (p pagination) IsDefined() bool {
	return p.isDefined
}
