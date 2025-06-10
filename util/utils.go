package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenStringId() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
