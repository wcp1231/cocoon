package mock

import (
	"os"
	"regexp"
)

type FieldMatcher interface {
	Match(field string) bool
}

type StringMatcher struct {
	Expect string
}

func (s *StringMatcher) Match(field string) bool {
	return s.Expect == field
}

type RegexMatcher struct {
	Regex *regexp.Regexp
}

func (r *RegexMatcher) Match(field string) bool {
	return r.Regex.MatchString(field)
}

func NewFieldMatcher(field *FieldMockConfig) FieldMatcher {
	if field.Equals != "" {
		return &StringMatcher{
			Expect: field.Equals,
		}
	}
	regex, err := regexp.Compile(field.Regex)
	if err != nil {
		// TODO
		os.Exit(1)
	}
	return &RegexMatcher{
		Regex: regex,
	}
}
