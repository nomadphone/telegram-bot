package phonenumbers

import (
	"regexp"
	"strings"
)

func NumbersOnly(phone string) string {
	return strings.Join(regexp.MustCompile(`\d`).FindAllString(phone, -1), "")
}
