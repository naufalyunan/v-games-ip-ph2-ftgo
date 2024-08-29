package models

type Response struct {
	Message string      `json:"message"`
	Note    string      `json:"note,omitempty"`
	Data    interface{} `json:"data"`
}

type ResponseDataRegister struct {
	ID           uint    `json:"id"`
	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	Deposit      float64 `json:"deposit,omitempty"`
	ReferralCode string  `json:"referral_code,omitempty"`
}
