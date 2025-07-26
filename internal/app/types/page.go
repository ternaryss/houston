package types

const (
	DefaultPage     int = 1
	DefaultPageSize int = 25
)

type Page[T any] struct {
	Current  int
	All      int
	Size     int
	Pageable int
	Content  []T
}

func NewPage[T any](pagination Pagination, content []T) Page[T] {
	return Page[T]{
		Current:  pagination.Page,
		All:      pagination.Pages,
		Size:     pagination.Size,
		Pageable: pagination.Pageable,
		Content:  content,
	}
}

func EmptyPage[T any]() Page[T] {
	return Page[T]{}
}

type Pagination struct {
	Page     int
	Pages    int
	Size     int
	Pageable int
	Limit    int
	Offset   int
}

func NewPagination(page, size, all int) Pagination {
	pages := 0

	if size > 0 && all > 0 {
		pages = all / size

		if all%size != 0 {
			pages = pages + 1
		}
	}

	return Pagination{
		Page:     page,
		Pages:    pages,
		Size:     size,
		Pageable: all,
		Limit:    size,
		Offset:   (page - 1) * size,
	}
}

func (p Pagination) IsValid() bool {
	if p.Pages == 0 || p.Page <= 0 || p.Size <= 0 || p.Page > p.Pages {
		return false
	}

	return true
}
