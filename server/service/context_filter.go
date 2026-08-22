package service

import (
	"regexp"
	"strings"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
)

const learnedContextMinimumOccurrences = 3

var simpleTokenPrefixPattern = regexp.MustCompile(`^[A-Za-z0-9_]+_$`)

// LearnedContextFilter contains content-only false-positive signatures learned
// from ignored results. A signature is enabled only when it occurs repeatedly
// in ignored results and never occurs in a confirmed result.
type LearnedContextFilter struct {
	patterns []learnedContextPattern
}

type learnedContextPattern struct {
	signature string
	pattern   *regexp.Regexp
}

// BuildLearnedContextFilter loads the historical results for one rule and
// learns conservative content signatures. Repository names and file paths are
// deliberately excluded so the filter remains useful when code moves.
func BuildLearnedContextFilter(keyword string) (*LearnedContextFilter, error) {
	filter := &LearnedContextFilter{}
	keyword = strings.TrimSpace(keyword)
	if global.GVA_DB == nil || !simpleTokenPrefixPattern.MatchString(keyword) {
		return filter, nil
	}

	var samples []model.SearchResult
	if err := global.GVA_DB.Select("status, text_matches_json, matches").
		Where("keyword = ? AND status IN ?", keyword, []int{global.IgnoredStatus, global.ConfirmedStatus}).
		Find(&samples).Error; err != nil {
		return nil, err
	}
	return learnContextFilter(samples, keyword), nil
}

func learnContextFilter(samples []model.SearchResult, keyword string) *LearnedContextFilter {
	ignored := make(map[string]int)
	confirmed := make(map[string]struct{})
	for _, sample := range samples {
		seen := make(map[string]struct{})
		for _, signature := range contextSignatures(SearchResultContent(sample), keyword) {
			if _, exists := seen[signature]; exists {
				continue
			}
			seen[signature] = struct{}{}
			if sample.Status == global.ConfirmedStatus {
				confirmed[signature] = struct{}{}
				continue
			}
			ignored[signature]++
		}
	}

	filter := &LearnedContextFilter{}
	for signature, count := range ignored {
		if count < learnedContextMinimumOccurrences {
			continue
		}
		if _, exists := confirmed[signature]; exists {
			continue
		}
		pattern, err := compileContextSignature(signature)
		if err != nil {
			continue
		}
		filter.patterns = append(filter.patterns, learnedContextPattern{
			signature: signature,
			pattern:   pattern,
		})
	}
	return filter
}

// ShouldFilter reports whether the result matches a learned false-positive
// signature. It intentionally returns only a boolean so raw snippets are not
// written to logs as part of the filtering decision.
func (f *LearnedContextFilter) ShouldFilter(result model.SearchResult) bool {
	if f == nil || len(f.patterns) == 0 {
		return false
	}
	evidence := SearchResultContent(result)
	for _, learned := range f.patterns {
		if learned.pattern.MatchString(evidence) {
			return true
		}
	}
	return false
}

func contextSignatures(content, keyword string) []string {
	if content == "" || keyword == "" {
		return nil
	}
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if signature, ok := normalizeContextLine(line, keyword); ok {
			result = append(result, signature)
		}
	}
	return result
}

func normalizeContextLine(line, keyword string) (string, bool) {
	line = strings.Join(strings.Fields(line), " ")
	if line == "" {
		return "", false
	}
	lowerLine := strings.ToLower(line)
	lowerKeyword := strings.ToLower(keyword)
	start := strings.Index(lowerLine, lowerKeyword)
	if start < 0 {
		return "", false
	}

	end := start + len(keyword)
	for end < len(line) && isTokenSuffixChar(line[end]) {
		end++
	}
	suffixLength := end - (start + len(keyword))
	marker := "<candidate-prefix>"
	switch {
	case suffixLength > 0 && suffixLength < 8:
		marker = "<short-candidate>"
	case suffixLength >= 8:
		marker = "<long-candidate>"
	}

	normalized := strings.ToLower(line[:start] + line[start:start+len(keyword)] + marker + line[end:])
	return normalized, true
}

func isTokenSuffixChar(char byte) bool {
	return char >= 'a' && char <= 'z' ||
		char >= 'A' && char <= 'Z' ||
		char >= '0' && char <= '9' || char == '_'
}

func compileContextSignature(signature string) (*regexp.Regexp, error) {
	pattern := regexp.QuoteMeta(signature)
	pattern = strings.ReplaceAll(pattern, " ", `\s+`)
	pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("<candidate-prefix>"), ``)
	pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("<short-candidate>"), `[A-Za-z0-9_]{1,7}`)
	pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("<long-candidate>"), `[A-Za-z0-9_]{8,}`)
	if strings.HasSuffix(signature, "<candidate-prefix>") || strings.HasSuffix(signature, "<short-candidate>") {
		pattern += `(?:[^A-Za-z0-9_]|$)`
	}
	return regexp.Compile(`(?i)` + pattern)
}
