package dto

type CreatePaymentRequest struct {
	Email  string  `json:"email" binding:"required,email"`
	Amount float64 `json:"amount" binding:"required"`
}

type UpdatePaymentRequest struct {
	Status string `json:"status" binding:"required"`
}

type PaymentResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}
