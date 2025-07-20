package handler

import (
	"net/http"
	"time"

	"payment-service/internal/delivery/http/dto"
	"payment-service/internal/domain"
	"payment-service/internal/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Create godoc
// @Summary Create a new payment
// @Description Create a new payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentRequest true "Create Payment Request"
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Router /payments [post]
func (h *PaymentHandler) Create(c echo.Context) error {
	var req dto.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	payment := &domain.Payment{
		Email:  req.Email,
		Amount: req.Amount,
		Status: "success", // default status
	}

	created, err := h.paymentService.Create(c.Request().Context(), payment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := dto.PaymentResponse{
		ID:        created.ID.Hex(),
		Email:     created.Email,
		Amount:    created.Amount,
		Status:    created.Status,
		CreatedAt: created.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusCreated, res)
}

// GetAll godoc
// @Summary Get all payments
// @Description Get all payments
// @Tags Payments
// @Produce json
// @Success 200 {array} dto.PaymentResponse
// @Failure 500 {object} map[string]string
// @Router /payments [get]
func (h *PaymentHandler) GetAll(c echo.Context) error {
	payments, err := h.paymentService.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []dto.PaymentResponse
	for _, p := range payments {
		response = append(response, dto.PaymentResponse{
			ID:        p.ID.Hex(),
			Email:     p.Email,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary Get payment by ID
// @Description Get payment by ID
// @Tags Payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	payment, err := h.paymentService.GetByID(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Payment not found")
	}

	res := dto.PaymentResponse{
		ID:        payment.ID.Hex(),
		Email:     payment.Email,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary Update payment by ID
// @Description Update payment status by ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param request body dto.UpdatePaymentRequest true "Update Payment Request"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} map[string]string
// @Router /payments/{id} [put]
func (h *PaymentHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var req dto.UpdatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	payment := &domain.Payment{
		Status: req.Status,
	}

	updated, err := h.paymentService.Update(c.Request().Context(), objectID, payment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res := dto.PaymentResponse{
		ID:        updated.ID.Hex(),
		Email:     updated.Email,
		Amount:    updated.Amount,
		Status:    updated.Status,
		CreatedAt: updated.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary Delete payment by ID
// @Description Delete payment by ID
// @Tags Payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/{id} [delete]
func (h *PaymentHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.paymentService.Delete(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
