package pagination

type OffsetBased struct {
	Page  int
	Limit int
}

func (o *OffsetBased) Offset() int {
	if o.Page > 1 {
		return o.Page * o.Limit
	}
	return 0
}

func SafeOffsetBased(page int, limit int) *OffsetBased {
	validPage := 1
	if page > 0 {
		validPage = page
	}
	validLimit := 1
	if limit > 0 {
		validLimit = limit
	}
	return &OffsetBased{
		Page:  validPage,
		Limit: validLimit,
	}
}

type OffsetBasedMetadata struct {
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	TotalPage int `json:"totalPage"`
}
