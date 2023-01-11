package request

type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type CasbinInReceive struct {
	AuthorityId string       `json:"authorityId"`
	CasbinInfos []CasbinInfo `json:"casbinInfos"`
}
