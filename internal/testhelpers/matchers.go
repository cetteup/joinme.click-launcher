package testhelpers

import (
	"strings"

	"github.com/golang/mock/gomock"
)

type stringContainsMatcher struct {
	substr string
}

func (m stringContainsMatcher) Matches(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return false
	}
	return strings.Contains(s, m.substr)
}

func (m stringContainsMatcher) String() string {
	return "contains substring " + m.substr
}

func StringContainsMatcher(substr string) gomock.Matcher {
	return stringContainsMatcher{substr: substr}
}
