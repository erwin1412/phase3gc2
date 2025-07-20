package handler

import (
	"net/http"

	"product-service/internal/delivery/http/dto"
	"product-service/internal/domain"
	"product-service/internal/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetAll godoc
// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) GetAll(c echo.Context) error {
	products, err := h.productService.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []dto.ProductResponse
	for _, p := range products {
		response = append(response, dto.ProductResponse{
			ID:    p.ID.Hex(),
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary Get product by ID
// @Description Get product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	product, err := h.productService.GetByID(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	response := dto.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}

	return c.JSON(http.StatusOK, response)
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) Create(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	product, err := h.productService.Create(
		c.Request().Context(),
		&domain.Product{
			Name:  req.Name,
			Price: req.Price,
			Stock: req.Stock,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := dto.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update godoc
// @Summary Update product by ID
// @Description Update product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductRequest true "Update Product Request"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var req dto.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	product, err := h.productService.Update(
		c.Request().Context(),
		objectID,
		&domain.Product{
			Name:  req.Name,
			Price: req.Price,
			Stock: req.Stock,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := dto.ProductResponse{
		ID:    product.ID.Hex(),
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Delete product by ID
// @Description Delete product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.productService.Delete(c.Request().Context(), objectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
