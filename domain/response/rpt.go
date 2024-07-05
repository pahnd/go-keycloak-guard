package response

import "encoding/json"

type MissingPermissionTicketResponse struct {
	Message          string `json:"message"`
	Code             int    `json:"code"`
	PermissionTicket string `json:"permissionTicket"`
}

func NewRPTResponse(permissionTicket string) *MissingPermissionTicketResponse {
	return &MissingPermissionTicketResponse{
		PermissionTicket: permissionTicket,
		Code:             401,
		Message:          "The request is missing the Requesting Party Token (RPT). Please obtain an RPT using the provided permission ticket.",
	}
}

func (m *MissingPermissionTicketResponse) ToJson() []byte {
	result, _ := json.Marshal(m)
	return result
}
