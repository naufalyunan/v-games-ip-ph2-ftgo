package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"v-games-ip-ph2-ftgo/models"

	xendit "github.com/xendit/xendit-go/v6"
	invoice "github.com/xendit/xendit-go/v6/invoice"
)

type XenditResponse struct {
	Invoice  *invoice.Invoice `json:"invoice"`
	Response *http.Response   `json:"response"`
	Error    error            `json:"error"`
}

type XenditService struct {
	Client  *xendit.APIClient
	APIKey  string
	BaseURL string
}

func NewXenditService() *XenditService {
	apiKey := os.Getenv("XENDIT_API_KEY")
	baseURL := "https://api.xendit.co"

	xnd := xendit.NewClient(apiKey)
	return &XenditService{
		Client:  xnd,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

func (s *XenditService) CreateInvoice(payment models.Payment, user *models.User) (*XenditResponse, error) {

	items := []invoice.InvoiceItem{}
	for _, item := range payment.Cart.CartItems {
		duration := item.CalculateDaysDifference()
		category := item.Game.Genre
		price := duration * item.Game.RentPrice
		el := invoice.InvoiceItem{
			Name:     item.Game.Name,
			Quantity: float32(item.Quantity),
			Price:    float32(price),
			Category: &category,
		}
		items = append(items, el)
	}

	cust := invoice.NewCustomerObject()
	cust.Email.Set(&user.Email)
	cust.GivenNames.Set(&user.FullName)
	custID := fmt.Sprintf("customer-id-%d", user.ID)
	cust.CustomerId.Set(&custID)

	description := ""

	if payment.CouponCode == "" {
		description = fmt.Sprintf("Dummy Invoice RMT006 for cart with ID %d. No coupon used.", payment.CartID)
	} else {
		description = fmt.Sprintf("Dummy Invoice RMT006 for cart with ID %d. Using coupon code %s", payment.CartID, payment.CouponCode)

	}

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest("dumm-external-id-RMT006", payment.PaymentPrice) // [REQUIRED] | CreateInvoiceRequest
	createInvoiceRequest.SetCurrency("IDR")
	createInvoiceRequest.SetItems(items)
	createInvoiceRequest.SetDescription(description)
	createInvoiceRequest.SetCustomer(*cust)
	resp, r, err := s.Client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return nil, err
	}

	result := XenditResponse{
		Invoice:  resp,
		Response: r,
		Error:    err,
	}

	return &result, nil
}
