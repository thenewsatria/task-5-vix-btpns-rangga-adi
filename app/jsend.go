package app

type JsendSuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type JsendFailResponse struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type JsendErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
