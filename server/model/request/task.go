package request

type SwitchTaskReq struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}
