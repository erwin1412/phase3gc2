package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"gateway-service/config"
	pb "gateway-service/internal/pb"

	"github.com/labstack/echo/v4"
)

type GatewayHandler struct {
	GRPC *config.GRPCClients
}

func NewGatewayHandler(grpcClients *config.GRPCClients) *GatewayHandler {
	return &GatewayHandler{GRPC: grpcClients}
}

// ✅ HANYA 1 FUNCTION proxyRequest, yang benar:
func proxyRequest(c echo.Context, targetBaseURL string) error {
	// STRIP prefix `/api` supaya /api/products -> /products
	path := strings.TrimPrefix(c.Request().RequestURI, "/api")

	req, err := http.NewRequest(
		c.Request().Method,
		targetBaseURL+path,
		c.Request().Body,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create proxy request"})
	}

	req.Header = c.Request().Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "Failed to connect to downstream service"})
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			c.Response().Header().Add(k, vv)
		}
	}
	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response(), resp.Body)
	return err
}

// === REST → REST ===
func (h *GatewayHandler) ProxyToProductService(c echo.Context) error {
	targetURL := os.Getenv("PRODUCT_URL") // TANPA `/api`
	return proxyRequest(c, targetURL)
}

func (h *GatewayHandler) ProxyToTransactionService(c echo.Context) error {
	targetURL := os.Getenv("TRANSACTION_URL")
	return proxyRequest(c, targetURL)
}

func (h *GatewayHandler) ProxyToPaymentService(c echo.Context) error {
	targetURL := os.Getenv("PAYMENT_URL")
	return proxyRequest(c, targetURL)
}

// === REST → GRPC ===
func (h *GatewayHandler) CreatePaymentGRPC(c echo.Context) error {
	var req pb.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	res, err := h.GRPC.PaymentClient.CreatePayment(context.Background(), &req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) GetAllPaymentsGRPC(c echo.Context) error {
	res, err := h.GRPC.PaymentClient.GetAllPayments(
		context.Background(),
		&pb.Empty{},
	)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) GetPaymentByIDGRPC(c echo.Context) error {
	id := c.Param("id") // asumsi route: /payments-grpc/:id
	res, err := h.GRPC.PaymentClient.GetPaymentByID(
		context.Background(),
		&pb.GetByIDRequest{Id: id},
	)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) UpdatePaymentGRPC(c echo.Context) error {
	var req pb.UpdatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	res, err := h.GRPC.PaymentClient.UpdatePayment(
		context.Background(),
		&req,
	)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) DeletePaymentGRPC(c echo.Context) error {
	id := c.Param("id")
	_, err := h.GRPC.PaymentClient.DeletePayment(
		context.Background(),
		&pb.GetByIDRequest{Id: id},
	)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Payment deleted"})
}

// === Proxy Helper ===
// func proxyRequest(c echo.Context, targetURL string) error {
// 	req, err := http.NewRequest(
// 		c.Request().Method,
// 		targetURL+c.Request().RequestURI,
// 		c.Request().Body,
// 	)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create proxy request"})
// 	}

// 	req.Header = c.Request().Header

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return c.JSON(http.StatusBadGateway, map[string]string{"error": "Failed to connect to downstream service"})
// 	}
// 	defer resp.Body.Close()

// 	for k, v := range resp.Header {
// 		for _, vv := range v {
// 			c.Response().Header().Add(k, vv)
// 		}
// 	}
// 	c.Response().WriteHeader(resp.StatusCode)
// 	_, err = io.Copy(c.Response(), resp.Body)
// 	return err
// }

func (h *GatewayHandler) RegisterGRPC(c echo.Context) error {
	var req pb.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	res, err := h.GRPC.AuthClient.Register(context.Background(), &req)
	if err != nil {
		return c.JSON(502, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, res)
}

func (h *GatewayHandler) LoginGRPC(c echo.Context) error {
	var req pb.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	res, err := h.GRPC.AuthClient.Login(context.Background(), &req)
	if err != nil {
		return c.JSON(502, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, res)
}
