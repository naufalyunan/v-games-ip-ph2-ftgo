package models

import invoice "github.com/xendit/xendit-go/v6/invoice"

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

type PaymentResponse struct {
	ID            uint             `json:"id"`
	PaymentPrice  float64          `json:"payment_price"`
	PaymentStatus string           `json:"payment_status"`
	PaymentMethod string           `json:"payment_method"`
	Provider      string           `json:"provider"`
	CouponCode    string           `json:"coupon_code,omitempty"`
	Invoice       *invoice.Invoice `json:"invoice"`

	//Assc
	Cart Cart `gorm:"foreignKey:CartID;references:ID" json:"cart"`
}
