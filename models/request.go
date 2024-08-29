package models

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateGameStockPayload struct {
	NewStock int `json:"stock"`
}

type UpdatePaymentStatusPayload struct {
	Status string `json:"status"`
}
