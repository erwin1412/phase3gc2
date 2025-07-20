package dto

type CreateTransactionRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	PaymentID string `json:"payment_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type UpdateTransactionRequest struct {
	Status string `json:"status" binding:"required"`
}

type TransactionResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	PaymentID string  `json:"payment_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}
