package dto

type OrderRequest struct {
	Items         []Item `json:"items"`
	UserID        string `json:"user_id"`
	Address       string `json:"address"`
	PaymentMethod string `json:"payment_method"`
}

type Item struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}
