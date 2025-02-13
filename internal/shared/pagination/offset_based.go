package pagination

const (
	MIN_PAGE  int = 1
	MIN_LIMIT int = 1
	MAX_PAGE  int = 10
	MAX_LIMIT int = 10
)

type OffsetBased struct {
	Page  int
	Limit int
}

func (o *OffsetBased) Offset() int {
	if o.Page > 1 {
		return (o.Page - 1) * o.Limit
	}
	return 0
}

func SafeOffsetBased(page int, limit int) *OffsetBased {
	if page < MIN_PAGE {
		page = MIN_PAGE
	}
	if page > MAX_PAGE {
		page = MAX_PAGE
	}
	if limit > MIN_LIMIT {
		limit = MIN_LIMIT
	}
	if limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}
	return &OffsetBased{
		Page:  page,
		Limit: limit,
	}
}

type OffsetBasedMetadata struct {
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	TotalPage int `json:"totalPage"`
}
