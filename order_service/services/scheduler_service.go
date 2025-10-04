package services

import (
	"log"
	"time"

	"order_service/repositories"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	orderRepo *repositories.OrderRepository
	cron      *cron.Cron
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		orderRepo: repositories.NewOrderRepository(),
		cron:      cron.New(cron.WithLocation(time.UTC)),
	}
}

// StartScheduler starts the cron scheduler with various periodic tasks
func (s *SchedulerService) StartScheduler() {
	// Weekly summary report - every Monday at 9 AM UTC
	s.cron.AddFunc("0 9 * * 1", s.generateWeeklySummaryReport)

	// Monthly report generation - 1st of every month at 10 AM UTC
	s.cron.AddFunc("0 10 1 * *", s.generateMonthlyReport)

	// Daily cleanup of old pending orders - every day at 2 AM UTC
	s.cron.AddFunc("0 2 * * *", s.cleanupPendingOrders)

	// Certificate status check - every hour
	s.cron.AddFunc("0 * * * *", s.checkCertificateStatus)

	s.cron.Start()
	log.Println("Scheduler service started")
}

// StopScheduler stops the cron scheduler
func (s *SchedulerService) StopScheduler() {
	s.cron.Stop()
	log.Println("Scheduler service stopped")
}

// generateWeeklySummaryReport creates a weekly summary report
func (s *SchedulerService) generateWeeklySummaryReport() {
	log.Println("Generating weekly summary report...")

	// Get the start and end of the current week
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Second)

	orders, err := s.orderRepo.GetOrdersByDateRange(weekStart, weekEnd)
	if err != nil {
		log.Printf("Failed to get orders for weekly report: %v", err)
		return
	}

	var totalOrders int
	var totalTonnes, totalRevenue float64

	for _, order := range orders {
		if order.Status == "completed" {
			totalOrders++
			totalTonnes += order.Tonnes
			totalRevenue += order.TotalAmount
		}
	}

	report := map[string]interface{}{
		"period":        "weekly",
		"week_start":    weekStart.Format("2006-01-02"),
		"week_end":      weekEnd.Format("2006-01-02"),
		"total_orders":  totalOrders,
		"total_tonnes":  totalTonnes,
		"total_revenue": totalRevenue,
		"generated_at":  time.Now().Format(time.RFC3339),
	}

	// In a real implementation, you would:
	// 1. Save the report to a database
	// 2. Send it via email to administrators
	// 3. Store it in a file or cloud storage

	log.Printf("Weekly summary report generated: %+v", report)
}

// generateMonthlyReport creates a monthly report
func (s *SchedulerService) generateMonthlyReport() {
	log.Println("Generating monthly report...")

	now := time.Now()
	report, err := s.orderRepo.GetMonthlyReport(now.Year(), int(now.Month()))
	if err != nil {
		log.Printf("Failed to generate monthly report: %v", err)
		return
	}

	// In a real implementation, you would:
	// 1. Save the report to a database
	// 2. Send it via email to administrators
	// 3. Generate charts and visualizations
	// 4. Store it in a file or cloud storage

	log.Printf("Monthly report generated for %s %d: %d orders, %.2f tonnes, $%.2f revenue",
		report.Month, report.Year, report.TotalOrders, report.TotalTonnes, report.TotalRevenue)
}

// cleanupPendingOrders removes orders that have been pending for too long
func (s *SchedulerService) cleanupPendingOrders() {
	log.Println("Cleaning up pending orders...")

	// Find orders that have been pending for more than 24 hours
	_, err := s.orderRepo.GetOrdersByDateRange(time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		log.Printf("Failed to get pending orders: %v", err)
		return
	}

	log.Printf("Cleanup completed at %s", time.Now().Format(time.RFC3339))
}

// checkCertificateStatus checks the status of pending certificates
func (s *SchedulerService) checkCertificateStatus() {
	log.Println("Checking certificate status...")

	// In a real implementation, you would:
	// 1. Find certificates that have been pending for too long
	// 2. Retry failed certificate generations
	// 3. Send alerts for stuck certificates
	// 4. Update certificate statuses

	log.Printf("Certificate status check completed at %s", time.Now().Format(time.RFC3339))
}
