package pagination

import (
	"errors"

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
	c := Collect(data)

	switch cursor.(type) {
	case *EmptyULIDCursor:
		if len(data) <= cursor.Limit() {
			next = NewEmptyULIDCursor(cursor.Limit())
		} else {
			var id *ulid.ULID

			if len(data) < cursor.Limit()+1 {
				id = retrieveID(c.Last())
			} else {
				id = retrieveID(c.At(c.LastIndex() - 2))
			}

			next, err = NewNextULIDCursor(id, cursor.Limit())
			if err != nil {
				return nil, err
			}
		}
		prev = NewEmptyULIDCursor(cursor.Limit())

		return &ULIDCursorMetadata[T]{
			Data: c.Slice(0, cursor.Limit()),
			Metadata: &CursorBasedMetadata{
				Next:     next.Cursor(),
				Previous: prev.Cursor(),
			},
		}, nil

	case *NextULIDCursor:
		if len(data) < cursor.Limit()+2 {
			next = NewEmptyULIDCursor(cursor.Limit())
		} else {
			next, err = NewNextULIDCursor(retrieveID(c.Penultimate()), cursor.Limit())
			if err != nil {
				return nil, err
			}
		}
		prev, err = NewPrevULIDCursor(cursor.ID(), cursor.Limit())
		if err != nil {
			return nil, err
		}
		return &ULIDCursorMetadata[T]{
			Data: c.Slice(1, cursor.Limit()+1),
			Metadata: &CursorBasedMetadata{
				Next:     next.Cursor(),
				Previous: prev.Cursor(),
			},
		}, nil

	case *PrevULIDCursor:
		if len(data) < cursor.Limit()+2 {
			prev = NewEmptyULIDCursor(cursor.Limit())
		} else {
			prev, err = NewPrevULIDCursor(retrieveID(c.Second()), cursor.Limit())
			if err != nil {
				return nil, err
			}
		}

		next, err = NewNextULIDCursor(cursor.ID(), cursor.Limit())
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
		return nil, errors.New("cursor cannot be nil")
	}
}

type Collection[T comparable] struct {
	data []T
	null T
}

func Collect[T comparable](data []T) *Collection[T] {
	return &Collection[T]{data: data}
}

func (d Collection[T]) Slice(i, j int) []T {
	if d.At(i) == d.null {
		return d.data
	}
	if d.At(j) == d.null {
		return d.data[i:]
	}
	return d.data[i:j]
}

func (d Collection[T]) At(i int) T {
	if d.LastIndex() < 0 {
		return d.null
	}
	if d.LastIndex() < i {
		return d.null
	}
	return d.data[i]
}

func (d Collection[T]) First() T {
	return d.At(0)
}

func (d Collection[T]) Second() T {
	return d.At(1)
}

func (d Collection[T]) Last() T {
	return d.At(d.LastIndex())
}

func (d Collection[T]) Penultimate() T {
	return d.At(d.LastIndex() - 1)
}

func (d Collection[T]) LastIndex() int {
	return len(d.data) - 1
}
