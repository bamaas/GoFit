package data

import (
	"log/slog"
	"database/sql"
	"fmt"

	// Import the sqlite driver for its side effects (database driver registration)
	_ "modernc.org/sqlite"
)

type CheckIn struct {
	UUID     string    `json:"uuid,omitempty"`
	Notes	 string    `json:"notes,omitempty"`
	Datetime int64 	   `json:"datetime"`
	Weight   float64   `json:"weight"`
	UserID   int64     `json:"-"`
}

type CheckInWithStats struct {
	CheckIn
	MovingAverage float64 `json:"moving_average"`
	WeightDifference float64 `json:"weight_difference"`
}

type CheckInModel struct {
	DB *sql.DB
	logger *slog.Logger
}

func (m *CheckInModel) Get(userID int64, UUID string) (CheckIn, error) {

	m.logger.Debug("Get check-in", "UUID", UUID)

	q := `
	SELECT uuid, datetime, weight, notes
	FROM checkins
	WHERE user_id=? AND uuid=?`

	r, err := m.DB.Query(q, userID, UUID)		// TODO: Use QueryRow
	if err != nil {
		return CheckIn{}, err
	}

	entries := []CheckIn{}
	for r.Next() {
		var e CheckIn
		scanErr := r.Scan(&e.UUID, &e.Datetime, &e.Weight, &e.Notes)
		if scanErr != nil {
			return CheckIn{}, scanErr
		}
		entries = append(entries, e)
	}

	// Verify the loop did not exit due to an error
	if err = r.Err(); err != nil {
		return CheckIn{}, err
	}

	if len(entries) == 0 {
		return CheckIn{}, ErrRecordNotFound
	}

	return entries[0], nil
}

func (m *CheckInModel) List(userID int64, filters Filters) ([]CheckInWithStats, Metadata, error) {

	m.logger.Debug("Get all the check-ins")

	q := `
	SELECT count(*) OVER(), uuid, datetime, weight, notes, 
	avg(weight) OVER (
		ORDER BY datetime DESC
		RANGE BETWEEN 0 PRECEDING
		AND 6 * 24 * 60 * 60 FOLLOWING
	) AS MovingAverageWindow7,
	IFNULL(weight - LAG(weight, 1) OVER (ORDER BY datetime), 0.0) as weightDifference
	FROM checkins
	WHERE user_id=? %s
	ORDER BY datetime DESC
	LIMIT ?
	OFFSET ?;
	`
	args := []any{
		userID,
		filters.limit(),
		filters.offset(),
	}

	// Add date filter if provided
	if !filters.StartTime.IsZero() && !filters.EndTime.IsZero() {
		format := "2006-01-02"
		dateFilter := fmt.Sprintf("AND (datetime > strftime('%%s', '%s') AND datetime < strftime('%%s', '%s'))", filters.StartTime.Format(format), filters.EndTime.AddDate(0,0,1).Format(format))
		q = fmt.Sprintf(q, dateFilter)
	} else {
		q = fmt.Sprintf(q, "")
	}

	// Execute the query
	r, err := m.DB.Query(q, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	totalRecords := 0
	entries := []CheckInWithStats{}
	for r.Next() {
		var e CheckInWithStats
		scanErr := r.Scan(&totalRecords, &e.UUID, &e.Datetime, &e.Weight, &e.Notes, &e.MovingAverage, &e.WeightDifference)
		if scanErr != nil {
			return nil, Metadata{}, scanErr
		}
		entries = append(entries, e)
	}

	// Verify the loop did not exit due to an error
	if err = r.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return entries, metadata, nil
}

func (m *CheckInModel) Insert(checkIn CheckIn) error {

	m.logger.Debug("Insert check-in", "check-in", checkIn)

	q := `
	INSERT INTO checkins
	(uuid, user_id, datetime, weight, notes)
	VALUES
	(?, ?, ?, ?, ?);
	`
	_, err := m.DB.Exec(q, checkIn.UUID, checkIn.UserID, checkIn.Datetime, checkIn.Weight, checkIn.Notes)
	if err != nil {
		return err
	}
	return nil
}

func (m *CheckInModel) Delete(userID int64, UUID string) error {

	m.logger.Debug("Deleting", "UUID", UUID)

	q := `
	DELETE FROM checkins
	WHERE
	user_id=? AND uuid=?
	`
	_, err := m.DB.Exec(q, userID, UUID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CheckInModel) Update(checkIn CheckIn) error {

	m.logger.Debug("Updating", "UUID", checkIn.UUID)

	q := `
	UPDATE checkins
	SET weight=?, datetime=?, notes=?
	WHERE user_id=? AND uuid=?
	`
	_, err := m.DB.Exec(q, checkIn.Weight, checkIn.Datetime, checkIn.Notes, checkIn.UserID, checkIn.UUID)
	if err != nil {
		return err
	}
	return nil

}
