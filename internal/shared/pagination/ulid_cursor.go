package pagination

import (
	"encoding/base64"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

const (
	MAX_CURSOR_LIMIT = 10
	MIN_CURSOR_LIMIT = 1
)

type ULIDCursor interface {
	ID() *ulid.ULID
	Cursor() Base64Cursor
	Limit() int
}

type ULIDCursorPayload struct {
	ID        *ulid.ULID `json:"id"`
	Direction string     `json:"direction"`
}

func ParseULIDCursor(cursor string, limit int) (ULIDCursor, error) {
	if cursor == "" {
		return &EmptyULIDCursor{limit: limit}, nil
	}
	jsonCursor, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	payload := new(ULIDCursorPayload)
	if err := json.Unmarshal(jsonCursor, payload); err != nil {
		return nil, err
	}
	switch payload.Direction {
	case "prev":
		return &PrevULIDCursor{id: payload.ID, cursor: &cursor, limit: limit}, nil
	default:
		return &NextULIDCursor{id: payload.ID, cursor: &cursor, limit: limit}, nil
	}
}
