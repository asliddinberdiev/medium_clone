package models

type ResponseStatus struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseList struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Limit   uint        `json:"limit"`
	Page    uint        `json:"page"`
	Data    interface{} `json:"data"`
}
