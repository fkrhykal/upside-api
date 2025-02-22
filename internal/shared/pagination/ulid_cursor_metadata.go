package pagination

import (
	"errors"

	c "github.com/fkrhykal/upside-api/internal/shared/collection"
	"github.com/oklog/ulid/v2"
)

type ULIDCursorMetadata[T comparable] struct {
	Data     []T
	Metadata *CursorBasedMetadata
}

func ProcessULIDCursorMetadata[T comparable](data []T, cursor ULIDCursor, retrieveID func(item T) *ulid.ULID) (*ULIDCursorMetadata[T], error) {
	var next ULIDCursor
	var prev ULIDCursor
	var err error

	switch cursor.(type) {
	case *EmptyULIDCursor:
		if len(data) <= cursor.Limit() {
			next = NewEmptyULIDCursor()
		} else {
			var id *ulid.ULID

			if len(data) < cursor.Limit()+1 {
				id = retrieveID(c.Last(data))
			} else {
				id = retrieveID(c.At(data, c.LastIndex(data)-2))
			}

			next, err = NewNextULIDCursor(id)
			if err != nil {
				return nil, err
			}
		}
		prev = NewEmptyULIDCursor()

		return &ULIDCursorMetadata[T]{
			Data: c.Slice(data, 0, cursor.Limit()),
			Metadata: &CursorBasedMetadata{
				Next:     next.Cursor(),
				Previous: prev.Cursor(),
			},
		}, nil

	case *NextULIDCursor:
		if len(data) < cursor.Limit()+2 {
			next = NewEmptyULIDCursor()
		} else {
			next, err = NewNextULIDCursor(retrieveID(c.Penultimate(data)))
			if err != nil {
				return nil, err
			}
		}
		prev, err = NewPrevULIDCursor(cursor.ID())
		if err != nil {
			return nil, err
		}
		return &ULIDCursorMetadata[T]{
			Data: c.Slice(data, 1, cursor.Limit()+1),
			Metadata: &CursorBasedMetadata{
				Next:     next.Cursor(),
				Previous: prev.Cursor(),
			},
		}, nil

	case *PrevULIDCursor:
		if len(data) < cursor.Limit()+2 {
			prev = NewEmptyULIDCursor()
		} else {
			prev, err = NewPrevULIDCursor(retrieveID(c.Second(data)))
			if err != nil {
				return nil, err
			}
		}

		next, err = NewNextULIDCursor(cursor.ID())
		if err != nil {
			return nil, err
		}
		var prevData []T
		if len(data) < cursor.Limit()+2 {
			prevData = data
		} else {
			prevData = data[2:]
		}
		return &ULIDCursorMetadata[T]{
			Data:     prevData,
			Metadata: &CursorBasedMetadata{Next: next.Cursor(), Previous: prev.Cursor()},
		}, nil

	default:
		return nil, errors.New("cursor nil")
	}
}
