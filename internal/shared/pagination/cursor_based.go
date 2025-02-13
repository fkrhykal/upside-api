package pagination

type Base64Cursor *string

type CursorBasedMetadata struct {
	Next     Base64Cursor `json:"next"`
	Previous Base64Cursor `json:"previous"`
}
