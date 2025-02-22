package pagination

import (
	"encoding/base64"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

type NextULIDCursor struct {
	id     *ulid.ULID
	cursor Base64Cursor
	limit  int
}

func (n *NextULIDCursor) ID() *ulid.ULID {
	return n.id
}

func (n *NextULIDCursor) Cursor() Base64Cursor {
	return n.cursor
}

func (n *NextULIDCursor) Limit() int {
	return n.limit
}

func NewNextULIDCursor(ID *ulid.ULID) (*NextULIDCursor, error) {
	payload := &ULIDCursorPayload{
		ID:        ID,
		Direction: "next",
	}
	jsonCursor, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	cursor := base64.RawURLEncoding.EncodeToString(jsonCursor)
	return &NextULIDCursor{
		id:     ID,
		cursor: &cursor,
	}, nil
}
