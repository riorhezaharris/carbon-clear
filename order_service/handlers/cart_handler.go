package handlers

import (
	"net/http"
	"strconv"

	"order_service/models"
	"order_service/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartRepo  *repositories.CartRepository
	validator *validator.Validate
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartRepo:  repositories.NewCartRepository(),
		validator: validator.New(),
	}
}

// AddToCart adds items to the user's cart
func (h *CartHandler) AddToCart(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req models.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cartItem := &models.CartItem{
		UserID:    uint(userID),
		ProjectID: req.ProjectID,
		Tonnes:    req.Tonnes,
	}

	if err := h.cartRepo.AddToCart(cartItem); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add item to cart"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Item added to cart successfully"})
}

// GetCart retrieves the user's cart
func (h *CartHandler) GetCart(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	cartItems, err := h.cartRepo.GetCartByUserID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve cart"})
	}

	var responses []models.CartItemResponse
	for _, item := range cartItems {
		responses = append(responses, models.CartItemResponse{
			ID:        item.ID,
			UserID:    item.UserID,
			ProjectID: item.ProjectID,
			Tonnes:    item.Tonnes,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, responses)
}

// UpdateCartItem updates the quantity of an item in the cart
func (h *CartHandler) UpdateCartItem(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	projectIDStr := c.Param("projectID")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	var req models.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.cartRepo.UpdateCartItem(uint(userID), uint(projectID), req.Tonnes); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update cart item"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cart item updated successfully"})
}

// RemoveFromCart removes an item from the cart
func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	projectIDStr := c.Param("projectID")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	if err := h.cartRepo.RemoveFromCart(uint(userID), uint(projectID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove item from cart"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}

// ClearCart clears all items from the user's cart
func (h *CartHandler) ClearCart(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := h.cartRepo.ClearCart(uint(userID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to clear cart"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cart cleared successfully"})
}
