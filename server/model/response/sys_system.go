package response

import "github.com/madneal/gshark/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
