package response

import "github.com/madneal/gshark/model"

type ExaCustomerResponse struct {
	Customer model.ExaCustomer `json:"customer"`
}
