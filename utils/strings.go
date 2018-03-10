package utils

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/kennygrant/sanitize"
	"gopkg.in/tomb.v2"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

func StringContainsNoLetters(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func HasSuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// FormatTitle validates, escapes and capitalizes a title sentence, while lowercasing stop words.
func FormatTitle(s string) (string, error) {

	if !utf8.ValidString(s) {
		return "", errors.Newf("Invalid utf8 character in title [%s]", s)
	}

	s = strings.ToLower(s)
	s = sanitize.HTML(s)

	s = strings.Replace(s, `\`, `/`, -1)

	stopWords := []string{"a", "about", "an", "are", "as", "at", "be", "by", "com", "for", "from", "how", "in", "is", "it", "of", "on", "or", "that", "the", "this", "to", "was", "what", "when", "where", "who", "will", "with", "the"}

	fields := strings.Fields(s)
	for i := range fields {

		var err error
		fields[i], err = strconv.Unquote(`"` + strings.Replace(fields[i], `"`, "'", -1) + `"`)
		if err != nil {
			return "", errors.Wrapf(err, "Problem unquoting title [%s]", s)
		}

		fields[i] = strings.Trim(fields[i], `"`)

		if i == 0 {
			fields[i] = strings.Title(fields[i])
			continue
		}
		if ContainsString(stopWords, fields[i]) {
			fields[i] = strings.ToLower(fields[i])
			continue
		}
		fields[i] = strings.Title(fields[i])
	}

	return strings.Join(fields, " "), nil

}

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

func WordInArrayIsASubstring(title string, potentialSubstrings []string) bool {
	for _, a := range potentialSubstrings {
		if strings.Contains(title, a) {
			return true
		}
	}
	return false
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
