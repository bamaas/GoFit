package data

import (
	"log/slog"
	"time"

	"database/sql"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type CheckIn struct {
	UUID     string    `json:"uuid,omitempty"`
	Datetime int64 	   `json:"datetime"`
	Weight   float64   `json:"weight"`
	Notes	 string    `json:"notes,omitempty"`
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

func (m *CheckInModel) InjectSampleData() error {

	var checkIns []CheckIn

	for i := 1; i <= 58; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		checkIn := CheckIn{
			UUID:     uuid.String(),
			Datetime: time.Now().AddDate(0, 0, -i).Unix(),
			Weight:   float64(i + 29),
			Notes:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam nulla sem.",
		}
		checkIns = append(checkIns, checkIn)
	}

	for _, c := range checkIns {
		err := m.Insert(c)
		if err != nil {
			return err
		}
	}
	return nil

}

func (m *CheckInModel) Get(UUID string) (CheckIn, error) {

	m.logger.Debug("Get check-in", "UUID", UUID)

	q := `
	SELECT uuid, datetime, weight, notes
	FROM checkins
	WHERE uuid=?`

	r, err := m.DB.Query(q, UUID)
	if err != nil {
		return CheckIn{}, err
	}

	entries := []CheckIn{}
	for r.Next() {
		var e CheckIn
		err := r.Scan(&e.UUID, &e.Datetime, &e.Weight, &e.Notes)
		if err != nil {
			return CheckIn{}, err
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

func (m *CheckInModel) List(filters Filters) ([]CheckInWithStats, Metadata, error) {

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
	ORDER BY datetime DESC
	LIMIT ?
	OFFSET ?;
	`
	r, err := m.DB.Query(q, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	totalRecords := 0
	entries := []CheckInWithStats{}
	for r.Next() {
		var e CheckInWithStats
		err := r.Scan(&totalRecords, &e.UUID, &e.Datetime, &e.Weight, &e.Notes, &e.MovingAverage, &e.WeightDifference)
		if err != nil {
			return nil, Metadata{}, err
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
	(uuid, datetime, weight, notes)
	VALUES
	(?, ?, ?, ?);
	`
	_, err := m.DB.Exec(q, checkIn.UUID, checkIn.Datetime, checkIn.Weight, checkIn.Notes)
	if err != nil {
		return err
	}
	return nil
}

func (m *CheckInModel) Delete(UUID string) error {

	m.logger.Debug("Deleting", "UUID", UUID)

	q := `
	DELETE FROM checkins
	WHERE
	uuid=?
	`
	_, err := m.DB.Exec(q, UUID)
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
	WHERE uuid=?
	`
	_, err := m.DB.Exec(q, checkIn.Weight, checkIn.Datetime, checkIn.Notes, checkIn.UUID)
	if err != nil {
		return err
	}
	return nil

}
