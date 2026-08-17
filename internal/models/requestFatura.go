package models

type RequestFatura struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
