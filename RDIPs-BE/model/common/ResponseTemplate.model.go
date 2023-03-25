package model

type ResponseTemplate struct {
	HttpCode int         `json:"httpCode"`
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Error    error
}

func (response *ResponseTemplate) SetHttpCode(code int) {
	response.HttpCode = code
}

func (response *ResponseTemplate) SetData(data interface{}) {
	response.Data = data
}

func (response *ResponseTemplate) SetMessage(message string) {
	response.Message = message
}
