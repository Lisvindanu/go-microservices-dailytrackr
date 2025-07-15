package handlers

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	"dailytrackr/activity-service/models"
	"dailytrackr/activity-service/services"
	"dailytrackr/shared/config"
	"dailytrackr/shared/constants"
	"dailytrackr/shared/dto"

	"github.com/gofiber/fiber/v2"
)

type ActivityHandlers struct {
	activityRepo *models.ActivityRepository
	photoService *services.PhotoService
	config       *config.Config
}

// NewActivityHandlers creates a new activity handlers instance
func NewActivityHandlers(db *sql.DB, cfg *config.Config) *ActivityHandlers {
	return &ActivityHandlers{
		activityRepo: models.NewActivityRepository(db),
		photoService: services.NewPhotoService(cfg),
		config:       cfg,
	}
}

// CreateActivity handles creating a new activity
func (h *ActivityHandlers) CreateActivity(c *fiber.Ctx) error {
	// Get user ID from middleware (will be set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	var req dto.CreateActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Parse start time
	startTime, err := time.Parse(constants.DateTimeFormat, req.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid start_time format. Use: 2006-01-02T15:04:05Z",
			"error":   err.Error(),
		})
	}

	// Create activity
	activity := &models.Activity{
		UserID:       userID.(int64),
		Title:        req.Title,
		StartTime:    startTime,
		DurationMins: req.DurationMins,
		Cost:         req.Cost,
		Note:         req.Note,
	}

	if err := h.activityRepo.Create(activity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create activity",
			"error":   err.Error(),
		})
	}

	// Convert to response DTO
	response := h.convertToActivityResponse(activity)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": constants.MsgActivityCreated,
		"data":    response,
	})
}

// GetActivities handles getting user activities with pagination
func (h *ActivityHandlers) GetActivities(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Check for date range filter
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var activities []models.Activity
	var total int
	var err error

	if startDateStr != "" && endDateStr != "" {
		// Parse dates
		startDate, err := time.Parse(constants.DateFormat, startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid start_date format. Use: 2006-01-02",
			})
		}

		endDate, err := time.Parse(constants.DateFormat, endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid end_date format. Use: 2006-01-02",
			})
		}

		// Set end date to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

		activities, err = h.activityRepo.GetActivitiesByDateRange(userID.(int64), startDate, endDate)
		total = len(activities)

		// Apply pagination to filtered results
		start := offset
		end := offset + limit
		if start > len(activities) {
			activities = []models.Activity{}
		} else {
			if end > len(activities) {
				end = len(activities)
			}
			activities = activities[start:end]
		}
	} else {
		// Get activities with pagination
		activities, total, err = h.activityRepo.GetByUserID(userID.(int64), limit, offset)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get activities",
			"error":   err.Error(),
		})
	}

	// Convert to response DTOs
	var responses []dto.ActivityResponse
	for _, activity := range activities {
		responses = append(responses, h.convertToActivityResponse(&activity))
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	listResponse := dto.ActivityListResponse{
		Activities: responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"message":     "Activities retrieved successfully",
		"data":        listResponse,
		"total_pages": totalPages,
	})
}

// GetActivity handles getting a single activity by ID
func (h *ActivityHandlers) GetActivity(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	activityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid activity ID",
		})
	}

	var activity models.Activity
	err = h.activityRepo.GetByID(activityID, userID.(int64), &activity)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrActivityNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get activity",
			"error":   err.Error(),
		})
	}

	response := h.convertToActivityResponse(&activity)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Activity retrieved successfully",
		"data":    response,
	})
}

// UpdateActivity handles updating an activity
func (h *ActivityHandlers) UpdateActivity(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	activityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid activity ID",
		})
	}

	var req dto.UpdateActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Get existing activity
	var activity models.Activity
	err = h.activityRepo.GetByID(activityID, userID.(int64), &activity)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrActivityNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get activity",
			"error":   err.Error(),
		})
	}

	// Update fields if provided
	if req.Title != "" {
		activity.Title = req.Title
	}
	if req.StartTime != "" {
		startTime, err := time.Parse(constants.DateTimeFormat, req.StartTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid start_time format. Use: 2006-01-02T15:04:05Z",
			})
		}
		activity.StartTime = startTime
	}
	if req.DurationMins > 0 {
		activity.DurationMins = req.DurationMins
	}
	if req.Cost != nil {
		activity.Cost = req.Cost
	}
	if req.Note != "" {
		activity.Note = req.Note
	}

	if err := h.activityRepo.Update(&activity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update activity",
			"error":   err.Error(),
		})
	}

	// Get updated activity
	err = h.activityRepo.GetByID(activityID, userID.(int64), &activity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get updated activity",
			"error":   err.Error(),
		})
	}

	response := h.convertToActivityResponse(&activity)

	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MsgActivityUpdated,
		"data":    response,
	})
}

// DeleteActivity handles deleting an activity
func (h *ActivityHandlers) DeleteActivity(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	activityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid activity ID",
		})
	}

	err = h.activityRepo.Delete(activityID, userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrActivityNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete activity",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MsgActivityDeleted,
	})
}

// UploadPhoto handles photo upload for an activity
func (h *ActivityHandlers) UploadPhoto(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": constants.ErrInvalidToken,
		})
	}

	activityID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid activity ID",
		})
	}

	// Check if activity exists
	var activity models.Activity
	err = h.activityRepo.GetByID(activityID, userID.(int64), &activity)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": constants.ErrActivityNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get activity",
			"error":   err.Error(),
		})
	}

	// Get file from form
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No photo file provided",
			"error":   err.Error(),
		})
	}

	// Upload photo using photo service
	photoURL, err := h.photoService.UploadPhoto(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to upload photo",
			"error":   err.Error(),
		})
	}

	// Update activity with photo URL
	err = h.activityRepo.UpdatePhotoURL(activityID, userID.(int64), photoURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update activity photo",
			"error":   err.Error(),
		})
	}

	// ✅ FIXED: Use common DTO PhotoUploadResponse
	response := dto.PhotoUploadResponse{
		URL:       photoURL,
		SecureURL: photoURL,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MsgPhotoUploaded,
		"data":    response,
	})
}

// convertToActivityResponse converts Activity model to ActivityResponse DTO
func (h *ActivityHandlers) convertToActivityResponse(activity *models.Activity) dto.ActivityResponse {
	return dto.ActivityResponse{
		ID:           activity.ID,
		UserID:       activity.UserID,
		Title:        activity.Title,
		StartTime:    activity.StartTime,
		DurationMins: activity.DurationMins,
		Cost:         activity.Cost,
		PhotoURL:     activity.PhotoURL,
		Note:         activity.Note,
		CreatedAt:    activity.CreatedAt,
		UpdatedAt:    activity.UpdatedAt,
	}
}
