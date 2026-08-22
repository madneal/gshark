package service

import (
	"encoding/json"
	"testing"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
)

func TestLearnContextFilterDistinguishesPrefixOnlyAndLongCandidates(t *testing.T) {
	samples := make([]model.SearchResult, 0, 4)
	for i := 0; i < 3; i++ {
		samples = append(samples, contextSample(t, global.IgnoredStatus, "gho_"))
	}
	samples = append(samples, contextSample(t, global.ConfirmedStatus, "gho_abcdefgh"))

	filter := learnContextFilter(samples, "gho_")
	if !filter.ShouldFilter(contextSample(t, global.UnhandledStatus, "gho_")) {
		t.Fatal("expected prefix-only candidate to be filtered")
	}
	if filter.ShouldFilter(contextSample(t, global.UnhandledStatus, "gho_abcdefgh")) {
		t.Fatal("did not expect confirmed long candidate shape to be filtered")
	}
}

func TestLearnContextFilterRejectsRepeatedShortCandidate(t *testing.T) {
	samples := make([]model.SearchResult, 0, 4)
	for i := 0; i < 3; i++ {
		samples = append(samples, contextSample(t, global.IgnoredStatus, `TOKEN=gho_xx`))
	}
	samples = append(samples, contextSample(t, global.ConfirmedStatus, `TOKEN=gho_abcdefgh`))

	filter := learnContextFilter(samples, "gho_")
	if !filter.ShouldFilter(contextSample(t, global.UnhandledStatus, `TOKEN=gho_123`)) {
		t.Fatal("expected repeated short candidate context to be filtered")
	}
	if filter.ShouldFilter(contextSample(t, global.UnhandledStatus, `TOKEN=gho_abcdefgh`)) {
		t.Fatal("did not expect long candidate context to be filtered")
	}
}

func TestLearnContextFilterDoesNotLearnConfirmedSignature(t *testing.T) {
	samples := make([]model.SearchResult, 0, 4)
	for i := 0; i < 3; i++ {
		samples = append(samples, contextSample(t, global.IgnoredStatus, `TOKEN=gho_xx`))
	}
	samples = append(samples, contextSample(t, global.ConfirmedStatus, `TOKEN=gho_yy`))

	filter := learnContextFilter(samples, "gho_")
	if filter.ShouldFilter(contextSample(t, global.UnhandledStatus, `TOKEN=gho_zz`)) {
		t.Fatal("did not expect a signature shared by confirmed results to be learned")
	}
}

func TestBuildLearnedContextFilterSkipsNonPrefixRules(t *testing.T) {
	previousDB := global.GVA_DB
	global.GVA_DB = nil
	t.Cleanup(func() { global.GVA_DB = previousDB })

	filter, err := BuildLearnedContextFilter("password")
	if err != nil {
		t.Fatal(err)
	}
	if filter.ShouldFilter(contextSample(t, global.UnhandledStatus, "password=example")) {
		t.Fatal("non-token-prefix rules should not build a learned filter")
	}
}

func contextSample(t *testing.T, status int, fragment string) model.SearchResult {
	t.Helper()
	encoded, err := json.Marshal([]model.TextMatch{{Fragment: &fragment}})
	if err != nil {
		t.Fatal(err)
	}
	return model.SearchResult{Status: status, TextMatchesJson: encoded}
}
