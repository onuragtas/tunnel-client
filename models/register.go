package models

type Register struct {
	Success bool
	Data    Data
}

type Data struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
