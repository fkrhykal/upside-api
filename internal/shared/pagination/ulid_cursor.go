package pagination

import (
	"encoding/base64"
	"encoding/json"
	"log"

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
		return LimitNextULIDCursor(limit), nil
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

func ULIDCursorMetadata(prev *PrevULIDCursor, next *NextULIDCursor) *CursorBasedMetadata {
	log.Printf("prev: %+v", prev)
	log.Printf("next: %+v", next)
	return &CursorBasedMetadata{
		Next:     next.Cursor(),
		Previous: prev.Cursor(),
	}
}
