package models

import (
	"database/sql"
	"time"
)

// Activity represents the activity model
type Activity struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Title        string    `json:"title" db:"title"`
	StartTime    time.Time `json:"start_time" db:"start_time"`
	DurationMins int       `json:"duration_mins" db:"duration_mins"`
	Cost         *int      `json:"cost,omitempty" db:"cost"`
	PhotoURL     string    `json:"photo_url,omitempty" db:"photo_url"`
	Note         string    `json:"note,omitempty" db:"note"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ActivityRepository handles database operations for activities
type ActivityRepository struct {
	db *sql.DB
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Create creates a new activity in the database
func (r *ActivityRepository) Create(activity *Activity) error {
	query := `
		INSERT INTO activities (user_id, title, start_time, duration_mins, cost, note) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		activity.UserID,
		activity.Title,
		activity.StartTime,
		activity.DurationMins,
		activity.Cost,
		activity.Note,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	activity.ID = id

	// Get the created activity to populate timestamps
	return r.GetByID(activity.ID, activity.UserID, activity)
}

// GetByID retrieves an activity by ID for a specific user
func (r *ActivityRepository) GetByID(id, userID int64, activity *Activity) error {
	query := `
		SELECT id, user_id, title, start_time, duration_mins, cost, photo_url, note, created_at, updated_at
		FROM activities 
		WHERE id = ? AND user_id = ?
	`

	var cost sql.NullInt64
	var photoURL, note sql.NullString

	err := r.db.QueryRow(query, id, userID).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Title,
		&activity.StartTime,
		&activity.DurationMins,
		&cost,
		&photoURL,
		&note,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Handle NULL values
	if cost.Valid {
		costInt := int(cost.Int64)
		activity.Cost = &costInt
	} else {
		activity.Cost = nil
	}

	if photoURL.Valid {
		activity.PhotoURL = photoURL.String
	} else {
		activity.PhotoURL = ""
	}

	if note.Valid {
		activity.Note = note.String
	} else {
		activity.Note = ""
	}

	return nil
}

// GetByUserID retrieves all activities for a user with pagination
func (r *ActivityRepository) GetByUserID(userID int64, limit, offset int) ([]Activity, int, error) {
	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM activities WHERE user_id = ?"
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get activities with pagination
	query := `
		SELECT id, user_id, title, start_time, duration_mins, cost, photo_url, note, created_at, updated_at
		FROM activities 
		WHERE user_id = ? 
		ORDER BY start_time DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		var cost sql.NullInt64
		var photoURL, note sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Title,
			&activity.StartTime,
			&activity.DurationMins,
			&cost,
			&photoURL,
			&note,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Handle NULL values
		if cost.Valid {
			costInt := int(cost.Int64)
			activity.Cost = &costInt
		} else {
			activity.Cost = nil
		}

		if photoURL.Valid {
			activity.PhotoURL = photoURL.String
		} else {
			activity.PhotoURL = ""
		}

		if note.Valid {
			activity.Note = note.String
		} else {
			activity.Note = ""
		}

		activities = append(activities, activity)
	}

	return activities, total, nil
}

// Update updates an activity
func (r *ActivityRepository) Update(activity *Activity) error {
	query := `
		UPDATE activities 
		SET title = ?, start_time = ?, duration_mins = ?, cost = ?, note = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query,
		activity.Title,
		activity.StartTime,
		activity.DurationMins,
		activity.Cost,
		activity.Note,
		activity.ID,
		activity.UserID,
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

// UpdatePhotoURL updates the photo URL for an activity
func (r *ActivityRepository) UpdatePhotoURL(id, userID int64, photoURL string) error {
	query := `
		UPDATE activities 
		SET photo_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.Exec(query, photoURL, id, userID)
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

// Delete deletes an activity
func (r *ActivityRepository) Delete(id, userID int64) error {
	query := "DELETE FROM activities WHERE id = ? AND user_id = ?"

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

// GetActivitiesByDateRange retrieves activities within a date range
func (r *ActivityRepository) GetActivitiesByDateRange(userID int64, startDate, endDate time.Time) ([]Activity, error) {
	query := `
		SELECT id, user_id, title, start_time, duration_mins, cost, photo_url, note, created_at, updated_at
		FROM activities 
		WHERE user_id = ? AND start_time >= ? AND start_time <= ?
		ORDER BY start_time DESC
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
		var photoURL, note sql.NullString

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Title,
			&activity.StartTime,
			&activity.DurationMins,
			&cost,
			&photoURL,
			&note,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle NULL values
		if cost.Valid {
			costInt := int(cost.Int64)
			activity.Cost = &costInt
		} else {
			activity.Cost = nil
		}

		if photoURL.Valid {
			activity.PhotoURL = photoURL.String
		} else {
			activity.PhotoURL = ""
		}

		if note.Valid {
			activity.Note = note.String
		} else {
			activity.Note = ""
		}

		activities = append(activities, activity)
	}

	return activities, nil
}
