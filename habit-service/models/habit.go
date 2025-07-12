package models

import (
	"database/sql"
	"time"
)

// Habit represents the habit model
type Habit struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	ReminderTime string    `json:"reminder_time,omitempty" db:"reminder_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// HabitLog represents the habit log model
type HabitLog struct {
	ID        int64     `json:"id" db:"id"`
	HabitID   int64     `json:"habit_id" db:"habit_id"`
	Date      time.Time `json:"date" db:"date"`
	Status    string    `json:"status" db:"status"` // DONE, SKIPPED, FAILED
	PhotoURL  string    `json:"photo_url,omitempty" db:"photo_url"`
	Note      string    `json:"note,omitempty" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// HabitRepository handles database operations for habits
type HabitRepository struct {
	db *sql.DB
}

// NewHabitRepository creates a new habit repository
func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

// Create creates a new habit in the database
func (r *HabitRepository) Create(habit *Habit) error {
	query := `
		INSERT INTO habits (user_id, title, start_date, end_date, reminder_time) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		habit.UserID,
		habit.Title,
		habit.StartDate,
		habit.EndDate,
		habit.ReminderTime,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	habit.ID = id
	return r.GetByID(habit.ID, habit.UserID, habit)
}

// GetByID retrieves a habit by ID for a specific user
func (r *HabitRepository) GetByID(id, userID int64, habit *Habit) error {
	query := `
		SELECT id, user_id, title, start_date, end_date, reminder_time, created_at, updated_at
		FROM habits 
		WHERE id = ? AND user_id = ?
	`

	var reminderTime sql.NullString
	err := r.db.QueryRow(query, id, userID).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Title,
		&habit.StartDate,
		&habit.EndDate,
		&reminderTime,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)

	if reminderTime.Valid {
		habit.ReminderTime = reminderTime.String
	}

	return err
}

// GetByUserID retrieves all habits for a user
func (r *HabitRepository) GetByUserID(userID int64) ([]Habit, error) {
	query := `
		SELECT id, user_id, title, start_date, end_date, reminder_time, created_at, updated_at
		FROM habits 
		WHERE user_id = ? 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		var reminderTime sql.NullString

		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Title,
			&habit.StartDate,
			&habit.EndDate,
			&reminderTime,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reminderTime.Valid {
			habit.ReminderTime = reminderTime.String
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

// Update updates a habit
func (r *HabitRepository) Update(habit *Habit) error {
	query := `
		UPDATE habits 
		SET title = ?, reminder_time = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		habit.Title,
		habit.ReminderTime,
		habit.ID,
		habit.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete deletes a habit
func (r *HabitRepository) Delete(id, userID int64) error {
	query := "DELETE FROM habits WHERE id = ? AND user_id = ?"

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetActiveHabits retrieves active habits for a user (habits that are currently running)
func (r *HabitRepository) GetActiveHabits(userID int64) ([]Habit, error) {
	query := `
		SELECT id, user_id, title, start_date, end_date, reminder_time, created_at, updated_at
		FROM habits 
		WHERE user_id = ? AND start_date <= CURDATE() AND end_date >= CURDATE()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		var reminderTime sql.NullString

		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Title,
			&habit.StartDate,
			&habit.EndDate,
			&reminderTime,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reminderTime.Valid {
			habit.ReminderTime = reminderTime.String
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

// HabitLogRepository handles database operations for habit logs
type HabitLogRepository struct {
	DB *sql.DB // Export field agar bisa diakses dari handlers
}

// NewHabitLogRepository creates a new habit log repository
func NewHabitLogRepository(db *sql.DB) *HabitLogRepository {
	return &HabitLogRepository{DB: db}
}

// Create creates a new habit log
func (r *HabitLogRepository) Create(log *HabitLog) error {
	query := `
		INSERT INTO habit_logs (habit_id, date, status, note) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE status = VALUES(status), note = VALUES(note), updated_at = CURRENT_TIMESTAMP
	`

	result, err := r.DB.Exec(query,
		log.HabitID,
		log.Date,
		log.Status,
		log.Note,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if id > 0 {
		log.ID = id
	}

	return r.GetByHabitAndDate(log.HabitID, log.Date, log)
}

// GetByHabitAndDate retrieves a habit log by habit ID and date
func (r *HabitLogRepository) GetByHabitAndDate(habitID int64, date time.Time, log *HabitLog) error {
	query := `
		SELECT id, habit_id, date, status, photo_url, note, created_at, updated_at
		FROM habit_logs 
		WHERE habit_id = ? AND date = ?
	`

	var photoURL, note sql.NullString
	err := r.DB.QueryRow(query, habitID, date.Format("2006-01-02")).Scan(
		&log.ID,
		&log.HabitID,
		&log.Date,
		&log.Status,
		&photoURL,
		&note,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if photoURL.Valid {
		log.PhotoURL = photoURL.String
	}
	if note.Valid {
		log.Note = note.String
	}

	return err
}

// GetByHabitID retrieves all logs for a habit
func (r *HabitLogRepository) GetByHabitID(habitID int64) ([]HabitLog, error) {
	query := `
		SELECT id, habit_id, date, status, photo_url, note, created_at, updated_at
		FROM habit_logs 
		WHERE habit_id = ? 
		ORDER BY date DESC
	`

	rows, err := r.DB.Query(query, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []HabitLog
	for rows.Next() {
		var log HabitLog
		var photoURL, note sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.HabitID,
			&log.Date,
			&log.Status,
			&photoURL,
			&note,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if photoURL.Valid {
			log.PhotoURL = photoURL.String
		}
		if note.Valid {
			log.Note = note.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// Update updates a habit log
func (r *HabitLogRepository) Update(log *HabitLog) error {
	query := `
		UPDATE habit_logs 
		SET status = ?, note = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.DB.Exec(query, log.Status, log.Note, log.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetLogByIDWithOwnership retrieves a habit log by ID and verifies user ownership
func (r *HabitLogRepository) GetLogByIDWithOwnership(logID, userID int64) (*HabitLog, error) {
	query := `
		SELECT hl.id, hl.habit_id, hl.date, hl.status, hl.photo_url, hl.note, hl.created_at, hl.updated_at
		FROM habit_logs hl
		JOIN habits h ON hl.habit_id = h.id
		WHERE hl.id = ? AND h.user_id = ?
	`

	var log HabitLog
	var photoURL, note sql.NullString
	err := r.DB.QueryRow(query, logID, userID).Scan(
		&log.ID, &log.HabitID, &log.Date, &log.Status, &photoURL, &note, &log.CreatedAt, &log.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if photoURL.Valid {
		log.PhotoURL = photoURL.String
	}
	if note.Valid {
		log.Note = note.String
	}

	return &log, nil
}

// GetStats calculates habit statistics
func (r *HabitLogRepository) GetStats(habitID int64) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_days,
			SUM(CASE WHEN status = 'DONE' THEN 1 ELSE 0 END) as completed_days,
			SUM(CASE WHEN status = 'SKIPPED' THEN 1 ELSE 0 END) as skipped_days,
			SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed_days
		FROM habit_logs 
		WHERE habit_id = ?
	`

	var totalDays, completedDays, skippedDays, failedDays int
	err := r.DB.QueryRow(query, habitID).Scan(&totalDays, &completedDays, &skippedDays, &failedDays)
	if err != nil {
		return nil, err
	}

	var successRate float64
	if totalDays > 0 {
		successRate = float64(completedDays) / float64(totalDays) * 100
	}

	// Calculate current streak
	currentStreak := r.calculateCurrentStreak(habitID)
	longestStreak := r.calculateLongestStreak(habitID)

	return map[string]interface{}{
		"total_days":     totalDays,
		"completed_days": completedDays,
		"skipped_days":   skippedDays,
		"failed_days":    failedDays,
		"success_rate":   successRate,
		"current_streak": currentStreak,
		"longest_streak": longestStreak,
	}, nil
}

// calculateCurrentStreak calculates the current streak of completed days
func (r *HabitLogRepository) calculateCurrentStreak(habitID int64) int {
	query := `
		SELECT status FROM habit_logs 
		WHERE habit_id = ? 
		ORDER BY date DESC 
		LIMIT 30
	`

	rows, err := r.DB.Query(query, habitID)
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

// calculateLongestStreak calculates the longest streak of completed days
func (r *HabitLogRepository) calculateLongestStreak(habitID int64) int {
	query := `
		SELECT status FROM habit_logs 
		WHERE habit_id = ? 
		ORDER BY date ASC
	`

	rows, err := r.DB.Query(query, habitID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	longestStreak := 0
	currentStreak := 0

	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			break
		}

		if status == "DONE" {
			currentStreak++
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}
		} else {
			currentStreak = 0
		}
	}

	return longestStreak
}
