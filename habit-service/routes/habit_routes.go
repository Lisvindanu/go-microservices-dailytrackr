package routes

import (
	"dailytrackr/habit-service/handlers"
	"dailytrackr/habit-service/middleware"

	"github.com/labstack/echo/v4"
)

// SetupHabitRoutes sets up all habit-related routes
func SetupHabitRoutes(e *echo.Echo, habitHandlers *handlers.HabitHandlers) {
	// API v1 routes with authentication
	api := e.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	// Habit routes
	habits := api.Group("/habits")
	{
		// CRUD operations
		habits.POST("", habitHandlers.CreateHabit)       // POST /api/v1/habits
		habits.GET("", habitHandlers.GetHabits)          // GET /api/v1/habits?active=true
		habits.GET("/:id", habitHandlers.GetHabit)       // GET /api/v1/habits/:id
		habits.PUT("/:id", habitHandlers.UpdateHabit)    // PUT /api/v1/habits/:id
		habits.DELETE("/:id", habitHandlers.DeleteHabit) // DELETE /api/v1/habits/:id

		// Habit logs
		habits.POST("/:id/logs", habitHandlers.CreateHabitLog) // POST /api/v1/habits/:id/logs
		habits.GET("/:id/logs", habitHandlers.GetHabitLogs)    // GET /api/v1/habits/:id/logs

		// Statistics
		habits.GET("/:id/stats", habitHandlers.GetHabitStats)       // GET /api/v1/habits/:id/stats
		habits.GET("/:id/complete", habitHandlers.GetHabitWithLogs) // GET /api/v1/habits/:id/complete
	}

	// Habit logs routes (for direct log management)
	habitLogs := api.Group("/habit-logs")
	{
		habitLogs.PUT("/:id", habitHandlers.UpdateHabitLog) // PUT /api/v1/habit-logs/:id
	}
}
