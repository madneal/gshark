package service

import (
	"testing"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gsharkModel "github.com/madneal/gshark/model"
	"github.com/stretchr/testify/require"
)

var _ persist.BatchAdapter = (*GormCasbinAdapter)(nil)

func TestSavePolicyLineMapsRulesToCasbinModel(t *testing.T) {
	line := savePolicyLine("p", []string{"100", "/api/users/:id", "GET", "domain", "eft", "extra"})

	require.Equal(t, &gsharkModel.CasbinModel{
		Ptype:       "p",
		AuthorityId: "100",
		Path:        "/api/users/:id",
		Method:      "GET",
		V3:          "domain",
		V4:          "eft",
		V5:          "extra",
	}, line)
}

func TestSavePolicyLineAllowsShortRules(t *testing.T) {
	line := savePolicyLine("p", []string{"100", "/api/users"})

	require.Equal(t, &gsharkModel.CasbinModel{
		Ptype:       "p",
		AuthorityId: "100",
		Path:        "/api/users",
	}, line)
}

func TestLoadPolicyLineLoadsStoredPolicy(t *testing.T) {
	model := casbinModel.NewModel()
	require.True(t, model.AddDef("p", "p", "sub, obj, act"))

	loadPolicyLine(gsharkModel.CasbinModel{
		Ptype:       "p",
		AuthorityId: "100",
		Path:        "/api/users/:id",
		Method:      "GET",
	}, model)

	require.Equal(t, [][]string{
		{"100", "/api/users/:id", "GET"},
	}, model["p"]["p"].Policy)
}
