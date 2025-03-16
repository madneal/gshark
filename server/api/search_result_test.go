package api

import (
	"encoding/json"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"testing"
)

func TestStartAITask(t *testing.T) {
	global.GVA_VP = initialize.Viper("../config.yaml")
	global.GVA_LOG = initialize.Zap()
	global.GVA_DB = initialize.Gorm()
	err, list := service.ListSearchResultByStatus(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, result := range list {
		var content string
		textMatches := make([]model.TextMatch, 0)
		if json.Valid(result.TextMatchesJson) {
			err = json.Unmarshal(result.TextMatchesJson, &textMatches)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, textMatch := range textMatches {
				content += *textMatch.Fragment + "\n"
			}
		} else {
			content = string(result.TextMatchesJson)
		}
		fmt.Println(content)
	}
}
