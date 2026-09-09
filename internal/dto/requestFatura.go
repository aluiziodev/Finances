package dto

type RequestFatura struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
