package pagination

import (
	"encoding/base64"
	"net/http"
)

const DefaultLimit = 50

type Cursor struct {
	After  string
	Before string
	Limit  int
}

func ParseCursor(request *http.Request) Cursor {
	query := request.URL.Query()

	cursor := Cursor{Limit: DefaultLimit}

	if after := query.Get("after"); after != "" {
		if decoded, err := base64.StdEncoding.DecodeString(after); err == nil {
			cursor.After = string(decoded)
		}
	}

	if before := query.Get("before"); before != "" {
		if decoded, err := base64.StdEncoding.DecodeString(before); err == nil {
			cursor.Before = string(decoded)
		}
	}

	return cursor
}

func EncodeCursor(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
}
