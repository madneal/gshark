package response

import "github.com/madneal/gshark/model"

type ExaFileResponse struct {
	File model.ExaFileUploadAndDownload `json:"file"`
}
