package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"dailytrackr/habit-service/models"
	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/dto"

	"github.com/labstack/echo/v4"
)

type HabitHandlers struct {
	habitRepo    *models.HabitRepository
	habitLogRepo *models.HabitLogRepository
	config       *config.Config
}

// NewHabitHandlers creates a new habit handlers instance
func NewHabitHandlers(db *sql.DB, cfg *config.Config) *HabitHandlers {
	return &HabitHandlers{
		habitRepo:    models.NewHabitRepository(db),
		habitLogRepo: models.NewHabitLogRepository(db),
		config:       cfg,
	}
}

// CreateHabit handles creating a new habit
func (h *HabitHandlers) CreateHabit(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	var req dto.CreateHabitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Parse dates
	startDate, err := time.Parse(constants.DateFormat, req.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid start_date format. Use: 2006-01-02",
			"error":   err.Error(),
		})
	}

	endDate, err := time.Parse(constants.DateFormat, req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid end_date format. Use: 2006-01-02",
			"error":   err.Error(),
		})
	}

	// Validate date range
	if endDate.Before(startDate) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "End date must be after start date",
		})
	}

	// Create habit
	habit := &models.Habit{
		UserID:       userID.(int64),
		Title:        req.Title,
		StartDate:    startDate,
		EndDate:      endDate,
		ReminderTime: req.ReminderTime,
	}

	if err := h.habitRepo.Create(habit); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create habit",
			"error":   err.Error(),
		})
	}

	response := convertToHabitResponse(habit)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": constants.MsgHabitCreated,
		"data":    response,
	})
}

// GetHabits handles getting user habits
func (h *HabitHandlers) GetHabits(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	// Check if only active habits are requested
	activeOnly := c.QueryParam("active") == "true"

	var habits []models.Habit
	var err error

	if activeOnly {
		habits, err = h.habitRepo.GetActiveHabits(userID.(int64))
	} else {
		habits, err = h.habitRepo.GetByUserID(userID.(int64))
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habits",
			"error":   err.Error(),
		})
	}

	var responses []dto.HabitResponse
	for _, habit := range habits {
		responses = append(responses, convertToHabitResponse(&habit))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Habits retrieved successfully",
		"data":    responses,
		"total":   len(responses),
	})
}

// GetHabit handles getting a single habit by ID
func (h *HabitHandlers) GetHabit(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit",
			"error":   err.Error(),
		})
	}

	response := convertToHabitResponse(&habit)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Habit retrieved successfully",
		"data":    response,
	})
}

// UpdateHabit handles updating a habit
func (h *HabitHandlers) UpdateHabit(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	var req dto.UpdateHabitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Get existing habit
	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit",
			"error":   err.Error(),
		})
	}

	// Update fields if provided
	if req.Title != "" {
		habit.Title = req.Title
	}
	if req.ReminderTime != "" {
		habit.ReminderTime = req.ReminderTime
	}

	if err := h.habitRepo.Update(&habit); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update habit",
			"error":   err.Error(),
		})
	}

	// Get updated habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get updated habit",
			"error":   err.Error(),
		})
	}

	response := convertToHabitResponse(&habit)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": constants.MsgHabitUpdated,
		"data":    response,
	})
}

// DeleteHabit handles deleting a habit
func (h *HabitHandlers) DeleteHabit(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	err = h.habitRepo.Delete(habitID, userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to delete habit",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": constants.MsgHabitDeleted,
	})
}

// CreateHabitLog handles creating/updating a habit log
func (h *HabitHandlers) CreateHabitLog(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	var req dto.CreateHabitLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Verify habit belongs to user
	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to verify habit",
			"error":   err.Error(),
		})
	}

	// Parse date
	date, err := time.Parse(constants.DateFormat, req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid date format. Use: 2006-01-02",
			"error":   err.Error(),
		})
	}

	// Create habit log
	log := &models.HabitLog{
		HabitID: habitID,
		Date:    date,
		Status:  req.Status,
		Note:    req.Note,
	}

	if err := h.habitLogRepo.Create(log); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create habit log",
			"error":   err.Error(),
		})
	}

	response := convertToHabitLogResponse(log)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": constants.MsgHabitLogCreated,
		"data":    response,
	})
}

// GetHabitLogs handles getting habit logs
func (h *HabitHandlers) GetHabitLogs(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	// Verify habit belongs to user
	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to verify habit",
			"error":   err.Error(),
		})
	}

	logs, err := h.habitLogRepo.GetByHabitID(habitID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit logs",
			"error":   err.Error(),
		})
	}

	var responses []dto.HabitLogResponse
	for _, log := range logs {
		responses = append(responses, convertToHabitLogResponse(&log))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Habit logs retrieved successfully",
		"data":    responses,
		"total":   len(responses),
	})
}

// UpdateHabitLog handles updating a habit log
func (h *HabitHandlers) UpdateHabitLog(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	logID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid log ID",
		})
	}

	var req dto.UpdateHabitLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Get existing log and verify ownership
	log, err := h.habitLogRepo.GetLogByIDWithOwnership(logID, userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitLogNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit log",
			"error":   err.Error(),
		})
	}

	// Update fields if provided
	if req.Status != "" {
		log.Status = req.Status
	}
	if req.Note != "" {
		log.Note = req.Note
	}

	if err := h.habitLogRepo.Update(log); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update habit log",
			"error":   err.Error(),
		})
	}

	response := convertToHabitLogResponse(log)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": constants.MsgHabitLogUpdated,
		"data":    response,
	})
}

// GetHabitStats handles getting habit statistics
func (h *HabitHandlers) GetHabitStats(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	// Verify habit belongs to user
	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to verify habit",
			"error":   err.Error(),
		})
	}

	stats, err := h.habitLogRepo.GetStats(habitID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit statistics",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Habit statistics retrieved successfully",
		"data":    stats,
	})
}

// GetHabitWithLogs handles getting habit with its logs and stats
func (h *HabitHandlers) GetHabitWithLogs(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	habitID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid habit ID",
		})
	}

	// Get habit
	var habit models.Habit
	err = h.habitRepo.GetByID(habitID, userID.(int64), &habit)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": constants.ErrHabitNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit",
			"error":   err.Error(),
		})
	}

	// Get logs
	logs, err := h.habitLogRepo.GetByHabitID(habitID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit logs",
			"error":   err.Error(),
		})
	}

	// Get stats
	stats, err := h.habitLogRepo.GetStats(habitID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get habit statistics",
			"error":   err.Error(),
		})
	}

	// Convert to response DTOs
	habitResponse := convertToHabitResponse(&habit)
	var logResponses []dto.HabitLogResponse
	for _, log := range logs {
		logResponses = append(logResponses, convertToHabitLogResponse(&log))
	}

	response := dto.HabitWithLogsResponse{
		Habit: habitResponse,
		Logs:  logResponses,
		Stats: dto.HabitStatsResponse{
			TotalDays:     stats["total_days"].(int),
			CompletedDays: stats["completed_days"].(int),
			SkippedDays:   stats["skipped_days"].(int),
			FailedDays:    stats["failed_days"].(int),
			SuccessRate:   stats["success_rate"].(float64),
			CurrentStreak: stats["current_streak"].(int),
			LongestStreak: stats["longest_streak"].(int),
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Habit with logs and stats retrieved successfully",
		"data":    response,
	})
}

// Helper functions
func convertToHabitResponse(habit *models.Habit) dto.HabitResponse {
	return dto.HabitResponse{
		ID:           habit.ID,
		UserID:       habit.UserID,
		Title:        habit.Title,
		StartDate:    habit.StartDate,
		EndDate:      habit.EndDate,
		ReminderTime: habit.ReminderTime,
		CreatedAt:    habit.CreatedAt,
		UpdatedAt:    habit.UpdatedAt,
	}
}

func convertToHabitLogResponse(log *models.HabitLog) dto.HabitLogResponse {
	return dto.HabitLogResponse{
		ID:        log.ID,
		HabitID:   log.HabitID,
		Date:      log.Date,
		Status:    log.Status,
		PhotoURL:  log.PhotoURL,
		Note:      log.Note,
		CreatedAt: log.CreatedAt,
		UpdatedAt: log.UpdatedAt,
	}
}
