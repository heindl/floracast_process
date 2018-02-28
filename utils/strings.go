package utils

import (
	"gopkg.in/tomb.v2"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

func CapitalizeString(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func ForEachStringToStrings(æ []string, callback func(string) ([]string, error)) ([]string, error) {
	res := []string{}
	lock := sync.Mutex{}
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _n := range æ {
			n := _n
			tmb.Go(func() error {
				s, err := callback(n)
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				res = AddStringToSet(res, s...)
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}
	return res, nil
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
