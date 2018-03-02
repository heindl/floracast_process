package utils

import (
	"github.com/dropbox/godropbox/errors"
	"io"
	"strings"
)

func ContainsError(a, b error) bool {
	if a == nil {
		return false
	}
	return strings.Contains(errors.GetMessage(a), errors.GetMessage(b))
}

func SafeClose(c io.Closer, err *error) {
	if closeErr := c.Close(); closeErr != nil && *err == nil {
		*err = closeErr
	}
}
