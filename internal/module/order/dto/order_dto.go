package dto

type Order struct {
	Id     int     `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}
