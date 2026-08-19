package api

import (
	"testing"

	"github.com/madneal/gshark/global"
	"go.uber.org/zap"
)

func TestConvertCsvIntoRules(t *testing.T) {
	if global.GVA_LOG == nil {
		global.GVA_LOG = zap.NewNop()
	}
	header := []string{"规则类型", "规则内容", "规则名称", "规则描述"}
	valid := []string{"github", "sec.com.cn in:file", "域名", "描述"}
	short := []string{"github", "sec.com.cn"}

	rules, skipped := convertCsvIntoRules([][]string{header, valid, short})
	if len(rules) != 1 {
		t.Fatalf("expected 1 converted rule, got %d", len(rules))
	}
	if skipped != 1 {
		t.Fatalf("expected 1 skipped row, got %d", skipped)
	}
	got := rules[0]
	if got.RuleType != "github" || got.Content != "sec.com.cn in:file" ||
		got.Name != "域名" || got.Desc != "描述" || !got.Status {
		t.Errorf("converted rule = %+v", got)
	}
}
