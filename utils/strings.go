package utils


import (
	"strings"
	"unicode/utf8"
	"unicode"
)

func CapitalizeString(s string) string {
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToUpper(r)) + s[n:]
}

func RemoveStringDuplicates(haystack []string) []string {

	res := []string{}

	for _, s := range haystack {
		if ContainsString(res, s) {
			continue
		}
		res = append(res, s)
	}

	return res
}

func IndexOfString(haystack []string, needle string) int {
	for i := range haystack {
		if needle == haystack[i] {
			return i
		}
	}
	return -1
}

func ContainsString(haystack []string, needle string) bool {
	return IndexOfString(haystack, needle) != -1
}

func IntersectsStrings(a []string, b []string) bool {
	for _, s := range a {
		if ContainsString(b, s) {
			return true
		}
	}
	return false
}


func StringsToLower(a ...string) []string {
	res := []string{}
	for _, s := range a {
		res = append(res, strings.ToLower(s))
	}
	return res
}

func AddStringToSet(haystack []string, needles ...string) []string {
	for _, needle := range needles {
		if needle == "" {
			continue
		}
		if !ContainsString(haystack, needle) {
			haystack = append(haystack, needle)
		}
	}
	return haystack
}
