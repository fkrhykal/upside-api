package pagination

import "github.com/oklog/ulid/v2"

type EmptyULIDCursor struct {
	limit int
}

func (e *EmptyULIDCursor) ID() *ulid.ULID {
	return nil
}

func (e *EmptyULIDCursor) Limit() int {
	return e.limit
}

func (e *EmptyULIDCursor) Cursor() Base64Cursor {
	return nil
}

func NewEmptyULIDCursor() *EmptyULIDCursor {
	return &EmptyULIDCursor{}
}
