package models

import (
	"database/sql"
	"time"
)

// StatRepository handles database operations for statistics
type StatRepository struct {
	db *sql.DB
}

// NewStatRepository creates a new stat repository
func NewStatRepository(db *sql.DB) *StatRepository {
	return &StatRepository{db: db}
}

// Dashboard statistics structure
type DashboardStats struct {
	TotalActivities int     `json:"total_activities"`
	TotalHours      float64 `json:"total_hours"`
	TotalExpenses   int     `json:"total_expenses"`
	ActiveHabits    int     `json:"active_habits"`
	CompletedHabits int     `json:"completed_habits"`
	AvgDailyHours   float64 `json:"avg_daily_hours"`
	StreakDays      int     `json:"streak_days"`
	ThisWeekHours   float64 `json:"this_week_hours"`
	LastWeekHours   float64 `json:"last_week_hours"`
	HoursGrowth     float64 `json:"hours_growth_percent"`
}

// Activity summary structure
type ActivitySummary struct {
	Period            string          `json:"period"`
	TotalActivities   int             `json:"total_activities"`
	TotalHours        float64         `json:"total_hours"`
	TotalExpenses     int             `json:"total_expenses"`
	AvgDuration       float64         `json:"avg_duration_mins"`
	MostProductiveDay string          `json:"most_productive_day"`
	TopCategories     []CategoryStats `json:"top_categories"`
}

type CategoryStats struct {
	Category   string  `json:"category"`
	Count      int     `json:"count"`
	TotalHours float64 `json:"total_hours"`
	Percentage float64 `json:"percentage"`
}

// Habit progress structures
type HabitProgressSummary struct {
	TotalHabits     int                   `json:"total_habits"`
	ActiveHabits    int                   `json:"active_habits"`
	CompletedHabits int                   `json:"completed_habits"`
	OverallSuccess  float64               `json:"overall_success_rate"`
	HabitDetails    []HabitProgressDetail `json:"habit_details"`
}

type HabitProgressDetail struct {
	HabitID       int64   `json:"habit_id"`
	Title         string  `json:"title"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	TotalDays     int     `json:"total_days"`
	CompletedDays int     `json:"completed_days"`
	SuccessRate   float64 `json:"success_rate"`
	CurrentStreak int     `json:"current_streak"`
	Status        string  `json:"status"` // active, completed, upcoming
}

// Chart data structures
type ChartData struct {
	Labels []string     `json:"labels"`
	Data   []ChartPoint `json:"data"`
}

type ChartPoint struct {
	Date       string  `json:"date"`
	Hours      float64 `json:"hours"`
	Activities int     `json:"activities"`
	Expenses   int     `json:"expenses"`
}

// Expense report structure
type ExpenseReport struct {
	Period             string            `json:"period"`
	TotalExpenses      int               `json:"total_expenses"`
	AverageDaily       float64           `json:"average_daily"`
	HighestDay         ExpenseDay        `json:"highest_day"`
	ExpensesByCategory []ExpenseCategory `json:"expenses_by_category"`
	DailyBreakdown     []ExpenseDay      `json:"daily_breakdown"`
}

type ExpenseDay struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
	Count  int    `json:"count"`
}

type ExpenseCategory struct {
	Category   string  `json:"category"`
	Amount     int     `json:"amount"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// GetDashboardStats retrieves dashboard statistics for a user
func (r *StatRepository) GetDashboardStats(userID int64) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Total activities
	err := r.db.QueryRow(`
		SELECT COUNT(*), 
		       COALESCE(SUM(duration_mins), 0) / 60.0,
		       COALESCE(SUM(cost), 0)
		FROM activities 
		WHERE user_id = ?
	`, userID).Scan(&stats.TotalActivities, &stats.TotalHours, &stats.TotalExpenses)
	if err != nil {
		return nil, err
	}

	// Active and completed habits
	err = r.db.QueryRow(`
		SELECT 
			SUM(CASE WHEN start_date <= CURDATE() AND end_date >= CURDATE() THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN end_date < CURDATE() THEN 1 ELSE 0 END) as completed
		FROM habits 
		WHERE user_id = ?
	`, userID).Scan(&stats.ActiveHabits, &stats.CompletedHabits)
	if err != nil {
		return nil, err
	}

	// Average daily hours (last 30 days)
	err = r.db.QueryRow(`
		SELECT COALESCE(AVG(daily_hours), 0)
		FROM (
			SELECT DATE(start_time) as activity_date, SUM(duration_mins) / 60.0 as daily_hours
			FROM activities 
			WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
			GROUP BY DATE(start_time)
		) daily_stats
	`, userID).Scan(&stats.AvgDailyHours)
	if err != nil {
		return nil, err
	}

	// This week hours
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(duration_mins), 0) / 60.0
		FROM activities 
		WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)
	`, userID).Scan(&stats.ThisWeekHours)
	if err != nil {
		return nil, err
	}

	// Last week hours
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(duration_mins), 0) / 60.0
		FROM activities 
		WHERE user_id = ? 
		AND start_time >= DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) + 7 DAY)
		AND start_time < DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)
	`, userID).Scan(&stats.LastWeekHours)
	if err != nil {
		return nil, err
	}

	// Calculate growth percentage
	if stats.LastWeekHours > 0 {
		stats.HoursGrowth = ((stats.ThisWeekHours - stats.LastWeekHours) / stats.LastWeekHours) * 100
	}

	// Calculate streak days (simplified - consecutive days with activities)
	rows, err := r.db.Query(`
		SELECT DATE(start_time) as activity_date
		FROM activities 
		WHERE user_id = ? 
		GROUP BY DATE(start_time)
		ORDER BY activity_date DESC
		LIMIT 30
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			continue
		}
		dates = append(dates, date)
	}

	// Calculate streak
	if len(dates) > 0 {
		today := time.Now().Truncate(24 * time.Hour)
		for i, date := range dates {
			expectedDate := today.AddDate(0, 0, -i)
			if date.Truncate(24 * time.Hour).Equal(expectedDate) {
				stats.StreakDays++
			} else {
				break
			}
		}
	}

	return stats, nil
}

// GetActivitySummary retrieves activity summary for a date range
func (r *StatRepository) GetActivitySummary(userID int64, startDate, endDate time.Time) (*ActivitySummary, error) {
	summary := &ActivitySummary{
		Period: startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
	}

	// Basic stats
	err := r.db.QueryRow(`
		SELECT COUNT(*), 
		       COALESCE(SUM(duration_mins), 0) / 60.0,
		       COALESCE(SUM(cost), 0),
		       COALESCE(AVG(duration_mins), 0)
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ?
	`, userID, startDate, endDate).Scan(
		&summary.TotalActivities,
		&summary.TotalHours,
		&summary.TotalExpenses,
		&summary.AvgDuration,
	)
	if err != nil {
		return nil, err
	}

	// Most productive day
	err = r.db.QueryRow(`
		SELECT DAYNAME(start_time), SUM(duration_mins) as total_mins
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ?
		GROUP BY DAYNAME(start_time), WEEKDAY(start_time)
		ORDER BY total_mins DESC
		LIMIT 1
	`, userID, startDate, endDate).Scan(&summary.MostProductiveDay, new(int))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return summary, nil
}

// GetAllHabitsProgress retrieves progress for all user habits
func (r *StatRepository) GetAllHabitsProgress(userID int64) (*HabitProgressSummary, error) {
	summary := &HabitProgressSummary{}

	// Get habit counts
	err := r.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN start_date <= CURDATE() AND end_date >= CURDATE() THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN end_date < CURDATE() THEN 1 ELSE 0 END) as completed
		FROM habits 
		WHERE user_id = ?
	`, userID).Scan(&summary.TotalHabits, &summary.ActiveHabits, &summary.CompletedHabits)
	if err != nil {
		return nil, err
	}

	// Get habit details
	rows, err := r.db.Query(`
		SELECT h.id, h.title, h.start_date, h.end_date,
		       COALESCE(stats.total_days, 0) as total_days,
		       COALESCE(stats.completed_days, 0) as completed_days,
		       CASE 
		           WHEN h.start_date > CURDATE() THEN 'upcoming'
		           WHEN h.end_date < CURDATE() THEN 'completed'
		           ELSE 'active'
		       END as status
		FROM habits h
		LEFT JOIN (
			SELECT habit_id,
			       COUNT(*) as total_days,
			       SUM(CASE WHEN status = 'DONE' THEN 1 ELSE 0 END) as completed_days
			FROM habit_logs
			GROUP BY habit_id
		) stats ON h.id = stats.habit_id
		WHERE h.user_id = ?
		ORDER BY h.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalSuccessRate := 0.0
	habitCount := 0

	for rows.Next() {
		var detail HabitProgressDetail
		var startDate, endDate time.Time

		err := rows.Scan(
			&detail.HabitID,
			&detail.Title,
			&startDate,
			&endDate,
			&detail.TotalDays,
			&detail.CompletedDays,
			&detail.Status,
		)
		if err != nil {
			continue
		}

		detail.StartDate = startDate.Format("2006-01-02")
		detail.EndDate = endDate.Format("2006-01-02")

		if detail.TotalDays > 0 {
			detail.SuccessRate = float64(detail.CompletedDays) / float64(detail.TotalDays) * 100
			totalSuccessRate += detail.SuccessRate
			habitCount++
		}

		// Calculate current streak (simplified)
		detail.CurrentStreak = r.calculateCurrentStreak(detail.HabitID)

		summary.HabitDetails = append(summary.HabitDetails, detail)
	}

	// Calculate overall success rate
	if habitCount > 0 {
		summary.OverallSuccess = totalSuccessRate / float64(habitCount)
	}

	return summary, nil
}

// GetSpecificHabitProgress retrieves progress for a specific habit
func (r *StatRepository) GetSpecificHabitProgress(userID, habitID int64) (*HabitProgressDetail, error) {
	detail := &HabitProgressDetail{}
	var startDate, endDate time.Time

	err := r.db.QueryRow(`
		SELECT h.id, h.title, h.start_date, h.end_date,
		       COALESCE(stats.total_days, 0) as total_days,
		       COALESCE(stats.completed_days, 0) as completed_days,
		       CASE 
		           WHEN h.start_date > CURDATE() THEN 'upcoming'
		           WHEN h.end_date < CURDATE() THEN 'completed'
		           ELSE 'active'
		       END as status
		FROM habits h
		LEFT JOIN (
			SELECT habit_id,
			       COUNT(*) as total_days,
			       SUM(CASE WHEN status = 'DONE' THEN 1 ELSE 0 END) as completed_days
			FROM habit_logs
			WHERE habit_id = ?
			GROUP BY habit_id
		) stats ON h.id = stats.habit_id
		WHERE h.user_id = ? AND h.id = ?
	`, habitID, userID, habitID).Scan(
		&detail.HabitID,
		&detail.Title,
		&startDate,
		&endDate,
		&detail.TotalDays,
		&detail.CompletedDays,
		&detail.Status,
	)
	if err != nil {
		return nil, err
	}

	detail.StartDate = startDate.Format("2006-01-02")
	detail.EndDate = endDate.Format("2006-01-02")

	if detail.TotalDays > 0 {
		detail.SuccessRate = float64(detail.CompletedDays) / float64(detail.TotalDays) * 100
	}

	detail.CurrentStreak = r.calculateCurrentStreak(detail.HabitID)

	return detail, nil
}

// GetActivityChartData retrieves chart data for activities
func (r *StatRepository) GetActivityChartData(userID int64, chartType string, period int) (*ChartData, error) {
	chart := &ChartData{}

	var query string
	var args []interface{}

	switch chartType {
	case "daily":
		query = `
			SELECT DATE(start_time) as chart_date,
			       SUM(duration_mins) / 60.0 as hours,
			       COUNT(*) as activities,
			       COALESCE(SUM(cost), 0) as expenses
			FROM activities 
			WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
			GROUP BY DATE(start_time)
			ORDER BY chart_date DESC
			LIMIT ?
		`
		args = []interface{}{userID, period, period}

	case "weekly":
		query = `
			SELECT DATE_SUB(DATE(start_time), INTERVAL WEEKDAY(start_time) DAY) as week_start,
			       SUM(duration_mins) / 60.0 as hours,
			       COUNT(*) as activities,
			       COALESCE(SUM(cost), 0) as expenses
			FROM activities 
			WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL ? WEEK)
			GROUP BY week_start
			ORDER BY week_start DESC
			LIMIT ?
		`
		args = []interface{}{userID, period, period}

	case "monthly":
		query = `
			SELECT DATE_FORMAT(start_time, '%Y-%m-01') as month_start,
			       SUM(duration_mins) / 60.0 as hours,
			       COUNT(*) as activities,
			       COALESCE(SUM(cost), 0) as expenses
			FROM activities 
			WHERE user_id = ? AND start_time >= DATE_SUB(CURDATE(), INTERVAL ? MONTH)
			GROUP BY month_start
			ORDER BY month_start DESC
			LIMIT ?
		`
		args = []interface{}{userID, period, period}

	default:
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var point ChartPoint
		var date time.Time

		err := rows.Scan(&date, &point.Hours, &point.Activities, &point.Expenses)
		if err != nil {
			continue
		}

		point.Date = date.Format("2006-01-02")
		chart.Data = append(chart.Data, point)
		chart.Labels = append(chart.Labels, point.Date)
	}

	return chart, nil
}

// GetExpenseReport retrieves expense report for a date range
func (r *StatRepository) GetExpenseReport(userID int64, startDate, endDate time.Time) (*ExpenseReport, error) {
	report := &ExpenseReport{
		Period: startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
	}

	// Total expenses
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(cost), 0), COUNT(*)
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ? AND cost IS NOT NULL
	`, userID, startDate, endDate).Scan(&report.TotalExpenses, new(int))
	if err != nil {
		return nil, err
	}

	// Calculate average daily
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days > 0 {
		report.AverageDaily = float64(report.TotalExpenses) / float64(days)
	}

	// Highest expense day
	err = r.db.QueryRow(`
		SELECT DATE(start_time), SUM(cost), COUNT(*)
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ? AND cost IS NOT NULL
		GROUP BY DATE(start_time)
		ORDER BY SUM(cost) DESC
		LIMIT 1
	`, userID, startDate, endDate).Scan(
		new(time.Time),
		&report.HighestDay.Amount,
		&report.HighestDay.Count,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Daily breakdown
	rows, err := r.db.Query(`
		SELECT DATE(start_time) as expense_date, 
		       COALESCE(SUM(cost), 0) as amount,
		       COUNT(*) as count
		FROM activities 
		WHERE user_id = ? AND start_time BETWEEN ? AND ? AND cost IS NOT NULL
		GROUP BY DATE(start_time)
		ORDER BY expense_date DESC
	`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day ExpenseDay
		var date time.Time

		err := rows.Scan(&date, &day.Amount, &day.Count)
		if err != nil {
			continue
		}

		day.Date = date.Format("2006-01-02")
		report.DailyBreakdown = append(report.DailyBreakdown, day)
	}

	return report, nil
}

// calculateCurrentStreak calculates the current streak for a habit
func (r *StatRepository) calculateCurrentStreak(habitID int64) int {
	rows, err := r.db.Query(`
		SELECT status FROM habit_logs 
		WHERE habit_id = ? 
		ORDER BY date DESC 
		LIMIT 30
	`, habitID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	streak := 0
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			break
		}

		if status == "DONE" {
			streak++
		} else {
			break
		}
	}

	return streak
}
