package service

import (
	"regexp"
	"strings"

	"github.com/madneal/gshark/model"
)

// ContextFilter applies an optional local expression to search evidence. It
// has no database or provider dependencies, so it can be reused by every
// search backend without changing the scan schedule.
type ContextFilter struct {
	pattern *regexp.Regexp
}

// NewContextFilter compiles a rule's local match expression. An empty pattern
// keeps the legacy behavior and accepts every provider result.
func NewContextFilter(expression string) (*ContextFilter, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return &ContextFilter{}, nil
	}
	pattern, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return &ContextFilter{pattern: pattern}, nil
}

// Matches reports whether the result contains evidence accepted by the local
// rule expression. SearchResultContent uses text-match fragments first and
// falls back to the legacy Matches field.
func (f *ContextFilter) Matches(result model.SearchResult) bool {
	if f == nil || f.pattern == nil {
		return true
	}
	return f.pattern.MatchString(SearchResultContent(result))
}
