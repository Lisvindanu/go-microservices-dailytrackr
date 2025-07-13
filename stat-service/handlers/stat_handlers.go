package handlers

import (
	"database/sql"
	"strconv"
	"time"

	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"
	"dailytrackr/stat-service/models"

	"github.com/gin-gonic/gin"
)

type StatHandlers struct {
	statRepo *models.StatRepository
	config   *config.Config
}

// NewStatHandlers creates a new stat handlers instance
func NewStatHandlers(db *sql.DB, cfg *config.Config) *StatHandlers {
	return &StatHandlers{
		statRepo: models.NewStatRepository(db),
		config:   cfg,
	}
}

// GetDashboard handles getting user dashboard statistics
func (h *StatHandlers) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	dashboard, err := h.statRepo.GetDashboardStats(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get dashboard statistics", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Dashboard statistics retrieved successfully", dashboard)
}

// GetActivitySummary handles getting activity summary statistics
func (h *StatHandlers) GetActivitySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(constants.DateFormat, startDateStr)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid start_date format. Use: 2006-01-02", err)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse(constants.DateFormat, endDateStr)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid end_date format. Use: 2006-01-02", err)
			return
		}
	} else {
		endDate = time.Now()
	}
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	summary, err := h.statRepo.GetActivitySummary(userID.(int64), startDate, endDate)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get activity summary", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Activity summary retrieved successfully", summary)
}

// GetHabitProgress handles getting habit progress statistics
func (h *StatHandlers) GetHabitProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	habitIDStr := c.Query("habit_id")
	var habitID int64
	var err error

	if habitIDStr != "" {
		habitID, err = strconv.ParseInt(habitIDStr, 10, 64)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid habit_id", err)
			return
		}
	}

	var progress interface{}
	if habitID > 0 {
		progress, err = h.statRepo.GetSpecificHabitProgress(userID.(int64), habitID)
	} else {
		progress, err = h.statRepo.GetAllHabitsProgress(userID.(int64))
	}

	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get habit progress", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Habit progress retrieved successfully", progress)
}

// GetActivityChart handles getting activity chart data
// file: stat_handlers.go

// GetActivityChart handles getting activity chart data
func (h *StatHandlers) GetActivityChart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// =================================================================
	// PERBAIKAN DI SINI
	// =================================================================

	// 1. Ambil parameter 'type', jika tidak ada, set default ke "daily"
	chartType := c.Query("type")
	if chartType == "" {
		chartType = "daily"
	}

	// 2. Ambil parameter 'period'
	periodStr := c.Query("period")
	var period int
	var err error

	if periodStr == "" {
		// 3. Jika tidak ada, set default ke 7
		period = 7
	} else {
		// 4. Jika ada, konversi ke integer
		period, err = strconv.Atoi(periodStr)
		if err != nil || period <= 0 {
			// Jika konversi gagal atau nilainya tidak valid, set default ke 7
			period = 7
		}
	}

	// =================================================================
	// AKHIR PERBAIKAN
	// =================================================================

	chartData, err := h.statRepo.GetActivityChartData(userID.(int64), chartType, period)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get chart data", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Activity chart data retrieved successfully", chartData)
}

// GetExpenseReport handles getting expense report
func (h *StatHandlers) GetExpenseReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(constants.DateFormat, startDateStr)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid start_date format. Use: 2006-01-02", err)
			return
		}
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse(constants.DateFormat, endDateStr)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid end_date format. Use: 2006-01-02", err)
			return
		}
	} else {
		endDate = time.Now()
	}
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	report, err := h.statRepo.GetExpenseReport(userID.(int64), startDate, endDate)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get expense report", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Expense report retrieved successfully", report)
}
