package utils

import (
	"github.com/dropbox/godropbox/errors"
	"strings"
)

func ContainsError(a, b error) bool {
	if a == nil {
		return false
	}
	return strings.Contains(errors.GetMessage(a), errors.GetMessage(b))
}
