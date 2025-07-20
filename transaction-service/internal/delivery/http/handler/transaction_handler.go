package handler

import (
	"net/http"
	"time"

	"transaction-service/internal/delivery/http/dto"
	"transaction-service/internal/domain"
	"transaction-service/internal/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetAll godoc
// @Summary Get all transactions
// @Description Get all transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Success 200 {array} dto.TransactionResponse
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h *TransactionHandler) GetAll(c echo.Context) error {
	transactions, err := h.transactionService.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []dto.TransactionResponse
	for _, t := range transactions {
		response = append(response, dto.TransactionResponse{
			ID:        t.ID.Hex(),
			ProductID: t.ProductID.Hex(),
			PaymentID: t.PaymentID.Hex(),
			Quantity:  t.Quantity,
			Total:     t.Total,
			Status:    t.Status,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary Get transaction by ID
// @Description Get transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	t, err := h.transactionService.GetByID(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Transaction not found")
	}

	res := dto.TransactionResponse{
		ID:        t.ID.Hex(),
		ProductID: t.ProductID.Hex(),
		PaymentID: t.PaymentID.Hex(),
		Quantity:  t.Quantity,
		Total:     t.Total,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, res)
}

// Create godoc
// @Summary Create transaction
// @Description Create new transaction with product & payment
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transaction Request"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
func (h *TransactionHandler) Create(c echo.Context) error {
	var req dto.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Product ID")
	}
	paymentID, err := primitive.ObjectIDFromHex(req.PaymentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Payment ID")
	}

	t := &domain.Transaction{
		ProductID: productID,
		PaymentID: paymentID, // ✅ Sekarang diisi!
		Quantity:  req.Quantity,

		// 👇 Email TIDAK di domain.Transaction! Jadi tetap lewat argumen service.
	}

	// ✅ PANGGIL SERVICE dengan email dari req.Email
	created, err := h.transactionService.Create(c.Request().Context(), t)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := dto.TransactionResponse{
		ID:        created.ID.Hex(),
		ProductID: created.ProductID.Hex(),
		PaymentID: created.PaymentID.Hex(),
		Quantity:  created.Quantity,
		Total:     created.Total,
		Status:    created.Status,
		CreatedAt: created.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusCreated, res)
}

// Update godoc
// @Summary Update transaction
// @Description Update transaction status
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body dto.UpdateTransactionRequest true "Update Transaction Request"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var req dto.UpdateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	t := &domain.Transaction{
		Status: req.Status,
	}

	updated, err := h.transactionService.Update(c.Request().Context(), objectID, t)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := dto.TransactionResponse{
		ID:        updated.ID.Hex(),
		ProductID: updated.ProductID.Hex(),
		PaymentID: updated.PaymentID.Hex(),
		Quantity:  updated.Quantity,
		Total:     updated.Total,
		Status:    updated.Status,
		CreatedAt: updated.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary Delete transaction
// @Description Delete transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.transactionService.Delete(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
