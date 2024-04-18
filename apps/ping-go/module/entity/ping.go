package entity

type PingRequest struct {
	Message string `json:"message"`
	Limit   int32  `json:"limit"`
}

type PingResponse struct {
	Message string `json:"message"`
}
