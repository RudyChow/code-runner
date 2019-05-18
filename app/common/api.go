package common

type ApiResponse struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
