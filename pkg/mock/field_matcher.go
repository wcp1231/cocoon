package mock

import (
	"regexp"
)

type FieldMatcher interface {
	Match(field string) bool
}

type StringMatcher struct {
	expect string
}

func (s *StringMatcher) Match(field string) bool {
	return s.expect == field
}

type RegexMatcher struct {
	expect *regexp.Regexp
}

func (r *RegexMatcher) Match(field string) bool {
	return r.expect.MatchString(field)
}