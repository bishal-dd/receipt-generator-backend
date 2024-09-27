package paginationUtil

import (
	"encoding/base64"
	"strconv"
	"strings"
)

func Offset(after *string) (int, error) {
	if after == nil {
		return 0, nil
	}

	decodedCursor, err := base64.StdEncoding.DecodeString(*after)
	if err != nil {
		return 0, err
	}

	cursorString := strings.TrimPrefix(string(decodedCursor), "cursor")
	start, err := strconv.Atoi(cursorString)
	if err != nil {
		return 0, err
	}

	return start, nil
}