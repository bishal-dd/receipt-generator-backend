package ids

import (
	"github.com/oklog/ulid/v2"
)

func UILD () string {
	id := ulid.Make()
	return id.String()
}