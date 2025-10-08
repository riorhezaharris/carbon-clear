package handlers

import (
	"net/http"
	"strconv"
	"time"

	"order_service/models"
	"order_service/repositories"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	orderRepo *repositories.OrderRepository
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		orderRepo: repositories.NewOrderRepository(),
	}
}

// GetMonthlyReport generates a monthly order report
// @Summary Get monthly report
// @Description Generate a monthly order report for admin (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param year query int false "Year"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} models.MonthlyReport "Monthly report"
// @Failure 400 {object} map[string]string "Invalid year or month"
// @Failure 500 {object} map[string]string "Failed to generate monthly report"
// @Router /api/v1/admin/reports/monthly [get]
func (h *AdminHandler) GetMonthlyReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	if yearStr == "" || monthStr == "" {
		// Default to current month if not specified
		now := time.Now()
		yearStr = strconv.Itoa(now.Year())
		monthStr = strconv.Itoa(int(now.Month()))
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid year"})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid month"})
	}

	if month < 1 || month > 12 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Month must be between 1 and 12"})
	}

	report, err := h.orderRepo.GetMonthlyReport(year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate monthly report"})
	}

	return c.JSON(http.StatusOK, report)
}

// GetOrdersByDateRange retrieves orders within a specific date range
// @Summary Get orders by date range
// @Description Retrieve orders within a specific date range (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} models.OrderResponse "List of orders"
// @Failure 400 {object} map[string]string "Invalid date format or missing parameters"
// @Failure 500 {object} map[string]string "Failed to retrieve orders"
// @Router /api/v1/admin/orders/date-range [get]
func (h *AdminHandler) GetOrdersByDateRange(c echo.Context) error {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date and end_date are required"})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format. Use YYYY-MM-DD"})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format. Use YYYY-MM-DD"})
	}

	orders, err := h.orderRepo.GetOrdersByDateRange(startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve orders"})
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		responses = append(responses, models.OrderResponse{
			ID:             order.ID,
			UserID:         order.UserID,
			ProjectID:      order.ProjectID,
			Tonnes:         order.Tonnes,
			PricePerTonne:  order.PricePerTonne,
			TotalAmount:    order.TotalAmount,
			Status:         order.Status,
			PaymentID:      order.PaymentID,
			CertificateURL: order.CertificateURL,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, responses)
}

// GetOrderStatistics provides general order statistics
// @Summary Get order statistics
// @Description Get overall order statistics with growth metrics (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security AdminAuth
// @Success 200 {object} map[string]interface{} "Order statistics"
// @Failure 500 {object} map[string]string "Failed to generate statistics"
// @Router /api/v1/admin/statistics [get]
func (h *AdminHandler) GetOrderStatistics(c echo.Context) error {
	// Get current month statistics
	now := time.Now()
	report, err := h.orderRepo.GetMonthlyReport(now.Year(), int(now.Month()))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate statistics"})
	}

	// Get last month statistics for comparison
	lastMonth := now.AddDate(0, -1, 0)
	lastMonthReport, err := h.orderRepo.GetMonthlyReport(lastMonth.Year(), int(lastMonth.Month()))
	if err != nil {
		// If we can't get last month's data, just return current month
		lastMonthReport = &models.MonthlyReport{}
	}

	statistics := map[string]interface{}{
		"current_month": map[string]interface{}{
			"month":         report.Month,
			"year":          report.Year,
			"total_orders":  report.TotalOrders,
			"total_tonnes":  report.TotalTonnes,
			"total_revenue": report.TotalRevenue,
		},
		"previous_month": map[string]interface{}{
			"month":         lastMonthReport.Month,
			"year":          lastMonthReport.Year,
			"total_orders":  lastMonthReport.TotalOrders,
			"total_tonnes":  lastMonthReport.TotalTonnes,
			"total_revenue": lastMonthReport.TotalRevenue,
		},
		"growth": map[string]interface{}{
			"orders_growth":  calculateGrowth(float64(report.TotalOrders), float64(lastMonthReport.TotalOrders)),
			"tonnes_growth":  calculateGrowth(report.TotalTonnes, lastMonthReport.TotalTonnes),
			"revenue_growth": calculateGrowth(report.TotalRevenue, lastMonthReport.TotalRevenue),
		},
	}

	return c.JSON(http.StatusOK, statistics)
}

func calculateGrowth(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	return ((current - previous) / previous) * 100.0
}
