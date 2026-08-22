package models

type RequestFatura struct {
	Description string `json:"description"`
	Bank        string `json:"bank"`
	Status      string `json:"status"`
}
