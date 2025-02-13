package pagination

import (
	"encoding/base64"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

type PrevULIDCursor struct {
	id     *ulid.ULID
	cursor Base64Cursor
	limit  int
}

func (p *PrevULIDCursor) ID() *ulid.ULID {
	return p.id
}

func (p *PrevULIDCursor) Cursor() Base64Cursor {
	return p.cursor
}

func (p *PrevULIDCursor) Limit() int {
	return p.limit
}

func LimitPrevULIDCursor(limit int) *PrevULIDCursor {
	return &PrevULIDCursor{
		limit: limit,
	}
}
func NewPrevULIDCursor(ID *ulid.ULID, limit int) (*PrevULIDCursor, error) {
	payload := &ULIDCursorPayload{
		ID:        ID,
		Direction: "prev",
	}
	jsonCursor, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	cursor := base64.RawURLEncoding.EncodeToString(jsonCursor)
	return &PrevULIDCursor{
		id:     ID,
		cursor: &cursor,
		limit:  limit,
	}, nil
}
