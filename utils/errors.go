package utils

import (
	"github.com/dropbox/godropbox/errors"
	"strings"
)

func ContainsError(a, b error) bool {
	return strings.Contains(errors.GetMessage(a), errors.GetMessage(b))
}
