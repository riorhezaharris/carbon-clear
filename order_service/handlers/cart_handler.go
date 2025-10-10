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
// @Summary Add item to cart
// @Description Add a project to the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param request body models.AddToCartRequest true "Cart item details"
// @Success 201 {object} map[string]string "Item added to cart successfully"
// @Failure 400 {object} map[string]string "Invalid user ID or request body"
// @Failure 500 {object} map[string]string "Failed to add item to cart"
// @Router /api/v1/cart/{userID}/items [post]
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
// @Summary Get user's cart
// @Description Retrieve all items in the user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} models.CartItemResponse "List of cart items"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Failed to retrieve cart"
// @Router /api/v1/cart/{userID} [get]
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
// @Summary Update cart item
// @Description Update the quantity of a specific item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param projectID path int true "Project ID"
// @Param request body models.UpdateCartItemRequest true "Updated cart item details"
// @Success 200 {object} map[string]string "Cart item updated successfully"
// @Failure 400 {object} map[string]string "Invalid user ID, project ID, or request body"
// @Failure 500 {object} map[string]string "Failed to update cart item"
// @Router /api/v1/cart/{userID}/items/{projectID} [put]
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

	var req models.UpdateCartItemRequest
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
// @Summary Remove item from cart
// @Description Remove a specific item from the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param projectID path int true "Project ID"
// @Success 200 {object} map[string]string "Item removed from cart successfully"
// @Failure 400 {object} map[string]string "Invalid user ID or project ID"
// @Failure 500 {object} map[string]string "Failed to remove item from cart"
// @Router /api/v1/cart/{userID}/items/{projectID} [delete]
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
// @Summary Clear cart
// @Description Clear all items from the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]string "Cart cleared successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Failed to clear cart"
// @Router /api/v1/cart/{userID} [delete]
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
