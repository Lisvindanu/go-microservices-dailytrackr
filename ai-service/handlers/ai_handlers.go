package handlers

import (
	"database/sql"
	"strconv"
	"time"

	"dailytrackr/ai-service/models"
	"dailytrackr/ai-service/services"
	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/utils"

	"github.com/gin-gonic/gin"
)

type AIHandlers struct {
	aiRepo    *models.AIRepository
	geminiSvc *services.GeminiService
	config    *config.Config
}

// NewAIHandlers creates a new AI handlers instance
func NewAIHandlers(db *sql.DB, cfg *config.Config) *AIHandlers {
	return &AIHandlers{
		aiRepo:    models.NewAIRepository(db),
		geminiSvc: services.NewGeminiService(cfg),
		config:    cfg,
	}
}

// GenerateDailySummary handles generating daily summary using AI
func (h *AIHandlers) GenerateDailySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Parse date (optional, defaults to today)
	dateStr := c.Query("date")
	var targetDate time.Time
	var err error

	if dateStr != "" {
		targetDate, err = time.Parse(constants.DateFormat, dateStr)
		if err != nil {
			utils.SendBadRequestResponse(c.Writer, "Invalid date format. Use: 2006-01-02", err)
			return
		}
	} else {
		targetDate = time.Now()
	}

	// Check if summary already exists for this date
	existingSummary, err := h.aiRepo.GetDailySummary(userID.(int64), targetDate)
	if err == nil && existingSummary != nil {
		utils.SendSuccessResponse(c.Writer, "Daily summary retrieved from cache", existingSummary)
		return
	}

	// Get user activities for the target date
	activities, err := h.aiRepo.GetUserActivitiesForDate(userID.(int64), targetDate)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get user activities", err)
		return
	}

	if len(activities) == 0 {
		utils.SendBadRequestResponse(c.Writer, "No activities found for the specified date", nil)
		return
	}

	// Generate AI summary
	summary, err := h.geminiSvc.GenerateDailySummary(activities)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to generate AI summary", err)
		return
	}

	// Save summary to database
	summaryRecord := &models.DailySummary{
		UserID:      userID.(int64),
		Date:        targetDate,
		SummaryText: summary,
		AIGenerated: true,
	}

	if err := h.aiRepo.SaveDailySummary(summaryRecord); err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to save summary", err)
		return
	}

	utils.SendSuccessResponse(c.Writer, "Daily summary generated successfully", summaryRecord)
}

// GenerateHabitRecommendation handles generating habit recommendations using AI
func (h *AIHandlers) GenerateHabitRecommendation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Parse days parameter (default 7 days)
	daysStr := c.Query("days")
	days := 7
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	// Get user activities from last N days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	activities, err := h.aiRepo.GetUserActivitiesForPeriod(userID.(int64), startDate, endDate)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get user activities", err)
		return
	}

	if len(activities) == 0 {
		utils.SendBadRequestResponse(c.Writer, "No activities found in the specified period", nil)
		return
	}

	// Get existing habits to avoid duplicates
	existingHabits, err := h.aiRepo.GetUserHabits(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get existing habits", err)
		return
	}

	// Generate AI recommendation
	recommendation, err := h.geminiSvc.GenerateHabitRecommendation(activities, existingHabits)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to generate habit recommendation", err)
		return
	}

	response := map[string]interface{}{
		"recommendation":   recommendation,
		"based_on_days":    days,
		"total_activities": len(activities),
		"existing_habits":  len(existingHabits),
		"analysis_period":  startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		"generated_at":     time.Now(),
	}

	utils.SendSuccessResponse(c.Writer, "Habit recommendation generated successfully", response)
}

// GetInsights handles getting AI insights for user
func (h *AIHandlers) GetInsights(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Get user data for insights
	insights, err := h.aiRepo.GetUserInsights(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get user insights", err)
		return
	}

	// Generate AI insights if user has enough data
	if insights.TotalActivities > 5 {
		aiInsights, err := h.geminiSvc.GenerateInsights(insights)
		if err == nil {
			insights.AIInsights = aiInsights
		}
	}

	utils.SendSuccessResponse(c.Writer, "User insights retrieved successfully", insights)
}

// AnalyzeActivities handles analyzing recent activities
func (h *AIHandlers) AnalyzeActivities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Parse period (default 7 days)
	daysStr := c.Query("days")
	days := 7
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 && parsedDays <= 30 {
			days = parsedDays
		}
	}

	// Get recent activities
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	activities, err := h.aiRepo.GetUserActivitiesForPeriod(userID.(int64), startDate, endDate)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get activities", err)
		return
	}

	if len(activities) == 0 {
		utils.SendBadRequestResponse(c.Writer, "No activities found for analysis", nil)
		return
	}

	// Generate analysis
	analysis, err := h.geminiSvc.AnalyzeActivities(activities, days)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to analyze activities", err)
		return
	}

	response := map[string]interface{}{
		"analysis":         analysis,
		"period_days":      days,
		"activities_count": len(activities),
		"period":           startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		"generated_at":     time.Now(),
	}

	utils.SendSuccessResponse(c.Writer, "Activity analysis completed successfully", response)
}

// GetProductivityTips handles getting AI-powered productivity tips
func (h *AIHandlers) GetProductivityTips(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendUnauthorizedResponse(c.Writer, constants.ErrInvalidToken)
		return
	}

	// Get user context for personalized tips
	userContext, err := h.aiRepo.GetUserContext(userID.(int64))
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to get user context", err)
		return
	}

	// Generate personalized tips
	tips, err := h.geminiSvc.GenerateProductivityTips(userContext)
	if err != nil {
		utils.SendInternalServerErrorResponse(c.Writer, "Failed to generate productivity tips", err)
		return
	}

	response := map[string]interface{}{
		"tips":         tips,
		"personalized": true,
		"based_on":     "User activity patterns and habits",
		"generated_at": time.Now(),
	}

	utils.SendSuccessResponse(c.Writer, "Productivity tips generated successfully", response)
}
