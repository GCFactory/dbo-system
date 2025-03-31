package models

type DefaultHttpResponse struct {
	Status int    `json:"status" validate:"required"`
	Info   string `json:"info" validate:"required"`
}

type GetAccountInfo struct {
	AccId string `json:"acc_id" validate:"required"`
}
