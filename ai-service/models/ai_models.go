package models

import (
	"database/sql"
	"time"
)

// AIRepository handles database operations for AI service
type AIRepository struct {
	db *sql.DB
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *sql.DB) *AIRepository {
	return &AIRepository{db: db}
}

// DailySummary represents daily summary model
type DailySummary struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Date        time.Time `json:"date"`
	SummaryText string    `json:"summary_text"`
	AIGenerated bool      `json:"ai_generated"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Activity represents activity model for AI analysis
type Activity struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	StartTime    time.Time `json:"start_time"`
	DurationMins int       `json:"duration_mins"`
	Cost         *int      `json:"cost,omitempty"`
	Note         string    `json:"note,omitempty"`
}

// Habit represents habit model for AI analysis
type Habit struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

// UserInsights represents comprehensive user insights
type UserInsights struct {
	UserID             int64     `json:"user_id"`
	TotalActivities    int       `json:"total_activities"`
	TotalHours         float64   `json:"total_hours"`
	TotalExpenses      int       `json:"total_expenses"`
	ActiveHabits       int       `json:"active_habits"`
	AvgDailyHours      float64   `json:"avg_daily_hours"`
	MostProductiveTime string    `json:"most_productive_time"`
	TopActivityType    string    `json:"top_activity_type"`
	SpendingPattern    string    `json:"spending_pattern"`
	AIInsights         string    `json:"ai_insights,omitempty"`
	LastUpdated        time.Time `json:"last_updated"`
}

// UserContext represents user context for personalized recommendations
type UserContext struct {
	UserID          int64   `json:"user_id"`
	Username        string  `json:"username"`
	TotalActivities int     `json:"total_activities"`
	TotalHabits     int     `json:"total_habits"`
	AvgDailyHours   float64 `json:"avg_daily_hours"`
	RecentPatterns  string  `json:"recent_patterns"`
}

// GetDailySummary retrieves daily summary for a user and date
func (r *AIRepository) GetDailySummary(userID int64, date time.Time) (*DailySummary, error) {
	summary := &DailySummary{}

	query := `
		SELECT id, user_id, date, summary_text, ai_generated, created_at, updated_at
		FROM daily_summary 
		WHERE user_id = ? AND date = ?
	`

	err := r.db.QueryRow(query, userID, date.Format("2006-01-02")).Scan(
		&summary.ID,
		&summary.UserID,
		&summary.Date,
		&summary.SummaryText,
		&summary.AIGenerated,
		&summary.CreatedAt,
		&summary.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return summary, nil
}

// SaveDailySummary saves daily summary to database
func (r *AIRepository) SaveDailySummary(summary *DailySummary) error {
	query := `
		INSERT INTO daily_summary (user_id, date, summary_text, ai_generated) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		summary_text = VALUES(summary_text), 
		ai_generated = VALUES(ai_generated),
		updated_at = CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(query,
		summary.UserID,
		summary.Date.Format("2006-01-02"),
		summary.SummaryText,
		summary.AIGenerated,
	)

	if err != nil {
		return err
	}

	if summary.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		summary.ID = id
	}

	return nil
}

// GetUserActivitiesForDate retrieves user activities for a specific date
func (r *AIRepository) GetUserActivitiesForDate(userID int64, date time.Time) ([]Activity, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24*time.Hour - time.Second)

	return r.GetUserActivitiesForPeriod(userID, startDate, endDate)
}

// GetUserActivitiesForPeriod retrieves user activities for a date range
func (r *AIRepository) GetUserActivitiesForPeriod(userID int64, startDate, endDate time.Time) ([]Activity, error) {
	query := `
		SELECT id, title, start_time, duration_mins, cost, note
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ?
		ORDER BY start_time ASC
	`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		var cost sql.NullInt64
		var note sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.Title,
			&activity.StartTime,
			&activity.DurationMins,
			&cost,
			&note,
		)
		if err != nil {
			continue
		}

		if cost.Valid {
			costInt := int(cost.Int64)
			activity.Cost = &costInt
		}

		if note.Valid {
			activity.Note = note.String
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// GetUserHabits retrieves user habits for recommendations
func (r *AIRepository) GetUserHabits(userID int64) ([]Habit, error) {
	query := `
		SELECT h.id, h.title,
		       CASE 
		           WHEN h.start_date > CURDATE() THEN 'upcoming'
		           WHEN h.end_date < CURDATE() THEN 'completed'
		           ELSE 'active'
		       END as status,
		       COALESCE(
		           (SELECT COUNT(*) * 100 / DATEDIFF(LEAST(h.end_date, CURDATE()), h.start_date)
		            FROM habit_logs hl 
		            WHERE hl.habit_id = h.id AND hl.status = 'DONE'), 0
		       ) as progress
		FROM habits h
		WHERE h.user_id = ?
		ORDER BY h.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit

		err := rows.Scan(
			&habit.ID,
			&habit.Title,
			&habit.Status,
			&habit.Progress,
		)
		if err != nil {
			continue
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

// GetUserInsights retrieves comprehensive user insights
func (r *AIRepository) GetUserInsights(userID int64) (*UserInsights, error) {
	insights := &UserInsights{
		UserID:      userID,
		LastUpdated: time.Now(),
	}

	// Basic activity stats
	err := r.db.QueryRow(`
		SELECT COUNT(*), 
		       COALESCE(SUM(duration_mins), 0) / 60.0,
		       COALESCE(SUM(cost), 0)
		FROM activities 
		WHERE user_id = ?
	`, userID).Scan(&insights.TotalActivities, &insights.TotalHours, &insights.TotalExpenses)

	if err != nil {
		return nil, err
	}

	// Active habits count
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM habits 
		WHERE user_id = ? AND start_date <= CURDATE() AND end_date >= CURDATE()
	`, userID).Scan(&insights.ActiveHabits)

	if err != nil {
		return nil, err
	}

	// Average daily hours (last 30 days)
	err = r.db.QueryRow(`
		SELECT COALESCE(AVG(daily_hours), 0)
		FROM (
			SELECT SUM(duration_mins) / 60.0 as daily_hours
			FROM activities 
			WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
			GROUP BY DATE(start_time)
		) daily_stats
	`, userID).Scan(&insights.AvgDailyHours)

	if err != nil {
		return nil, err
	}

	// Most productive time
	err = r.db.QueryRow(`
		SELECT HOUR(start_time) as hour, COUNT(*) as count
		FROM activities 
		WHERE user_id = ? 
		GROUP BY HOUR(start_time)
		ORDER BY count DESC
		LIMIT 1
	`, userID).Scan(new(int), new(int))

	if err == nil {
		// Could enhance this to return actual hour
		insights.MostProductiveTime = "Morning" // Simplified
	}

	// Top activity pattern (simplified)
	insights.TopActivityType = "Learning & Development"
	insights.SpendingPattern = "Moderate spender"

	return insights, nil
}

// GetUserContext retrieves user context for personalized recommendations
func (r *AIRepository) GetUserContext(userID int64) (*UserContext, error) {
	context := &UserContext{
		UserID: userID,
	}

	// Get username
	err := r.db.QueryRow(`
		SELECT username FROM users WHERE id = ?
	`, userID).Scan(&context.Username)

	if err != nil {
		return nil, err
	}

	// Get activity and habit counts
	err = r.db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM activities WHERE user_id = ?) as total_activities,
			(SELECT COUNT(*) FROM habits WHERE user_id = ?) as total_habits,
			COALESCE((SELECT AVG(daily_hours) FROM (
				SELECT SUM(duration_mins) / 60.0 as daily_hours
				FROM activities 
				WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
				GROUP BY DATE(start_time)
			) recent_daily), 0) as avg_daily_hours
	`, userID, userID, userID).Scan(
		&context.TotalActivities,
		&context.TotalHabits,
		&context.AvgDailyHours,
	)

	if err != nil {
		return nil, err
	}

	// Generate recent patterns summary
	if context.TotalActivities > 0 {
		if context.AvgDailyHours >= 4 {
			context.RecentPatterns = "High productivity, consistent daily activities"
		} else if context.AvgDailyHours >= 2 {
			context.RecentPatterns = "Moderate activity levels, room for improvement"
		} else {
			context.RecentPatterns = "Low activity levels, needs motivation boost"
		}
	} else {
		context.RecentPatterns = "New user, no patterns established yet"
	}

	return context, nil
}
