package service

import (
	"encoding/json"
	"testing"

	"github.com/madneal/gshark/model"
)

func TestContextFilterMatchesEvidence(t *testing.T) {
	filter, err := NewContextFilter(`(^|[^A-Za-z0-9_])gho_[A-Za-z0-9_]{4,}([^A-Za-z0-9_]|$)`)
	if err != nil {
		t.Fatal(err)
	}

	if !filter.Matches(contextFilterResult(t, `TOKEN="gho_abcdef"`)) {
		t.Fatal("expected a token-shaped fragment to match")
	}
	if filter.Matches(contextFilterResult(t, `TOKEN="gho_"`)) {
		t.Fatal("did not expect a prefix-only fragment to match")
	}
}

func TestContextFilterUsesLegacyMatchesFallback(t *testing.T) {
	filter, err := NewContextFilter(`ghp_[A-Za-z0-9_]{4,}`)
	if err != nil {
		t.Fatal(err)
	}
	if !filter.Matches(model.SearchResult{Matches: "ghp_abcdef"}) {
		t.Fatal("expected the legacy matches field to be checked")
	}
}

func TestContextFilterEmptyPatternPreservesBehavior(t *testing.T) {
	filter, err := NewContextFilter("")
	if err != nil {
		t.Fatal(err)
	}
	if !filter.Matches(contextFilterResult(t, "unrelated evidence")) {
		t.Fatal("empty expressions should accept legacy results")
	}
}

func TestContextFilterRejectsInvalidPattern(t *testing.T) {
	if _, err := NewContextFilter("("); err == nil {
		t.Fatal("expected invalid expressions to fail during rule setup")
	}
}

func contextFilterResult(t *testing.T, fragment string) model.SearchResult {
	t.Helper()
	encoded, err := json.Marshal([]model.TextMatch{{Fragment: &fragment}})
	if err != nil {
		t.Fatal(err)
	}
	return model.SearchResult{TextMatchesJson: encoded}
}
