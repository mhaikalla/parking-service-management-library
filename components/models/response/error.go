package response

type ErrorResponse struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type BaseMessageResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
