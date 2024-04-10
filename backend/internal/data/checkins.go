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
	Datetime time.Time `json:"datetime"`
	Weight   float64   `json:"weight"`
}

type CheckInModel struct {
	DB *sql.DB
	logger *slog.Logger
}

func (m *CheckInModel) InjectSampleData() error {

	var checkIns []CheckIn

	for i := 1; i <= 30; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		checkIn := CheckIn{
			UUID:     uuid.String(),
			Datetime: time.Now().AddDate(0, 0, -i),
			Weight:   float64(i + 29),
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

func parseRowsToEntries(r *sql.Rows) ([]CheckIn, error) {

	// Parse db sql rows
	type dbRow struct {
		UUID     string
		Datetime int64
		Weight   float64
	}
	dbData := []dbRow{}
	for r.Next() {
		var dbr dbRow
		err := r.Scan(&dbr.UUID, &dbr.Datetime, &dbr.Weight)
		if err != nil {
			return []CheckIn{}, err
		}
		dbData = append(dbData, dbr)
	}

	// Parse db data to CheckIn's/entries
	entries := []CheckIn{}
	for i := range dbData {
		entries = append(entries, CheckIn{
			UUID:     dbData[i].UUID,
			Datetime: time.Unix(dbData[i].Datetime, 0),
			Weight:   dbData[i].Weight})
	}

	return entries, nil
}

func (m *CheckInModel) Get(UUID string) (CheckIn, error) {

	m.logger.Debug("Get entry", "UUID", UUID)

	q := `
	SELECT uuid, datetime, weight
	FROM entries
	WHERE uuid=?`

	r, err := m.DB.Query(q, UUID)
	if err != nil {
		return CheckIn{}, err
	}

	entries, err := parseRowsToEntries(r)
	if err != nil {
		return CheckIn{}, err
	}

	return entries[0], nil
}

func (m *CheckInModel) List(filters Filters) ([]CheckIn, error) {

	m.logger.Debug("Get all the entries")

	q := `
	SELECT uuid, datetime, weight
	FROM entries
	ORDER BY datetime DESC
	LIMIT ?
	OFFSET ?;
	`
	r, err := m.DB.Query(q, filters.limit(), filters.offset())
	if err != nil {
		return nil, err
	}

	return parseRowsToEntries(r)
}

func (m *CheckInModel) Insert(checkIn CheckIn) error {

	m.logger.Debug("Insert check-in", "check-in", checkIn)

	q := `
	INSERT INTO entries
	(uuid, datetime, weight)
	VALUES
	(?, ?, ?);
	`
	_, err := m.DB.Exec(q, checkIn.UUID, checkIn.Datetime.Unix(), checkIn.Weight)
	if err != nil {
		return err
	}
	return nil
}

func (m *CheckInModel) Delete(UUID string) error {

	m.logger.Debug("Deleting", "UUID", UUID)

	q := `
	DELETE FROM entries
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
	UPDATE entries
	SET weight=?, datetime=?
	WHERE uuid=?
	`
	_, err := m.DB.Exec(q, checkIn.Weight, checkIn.Datetime.Unix(), checkIn.UUID)
	if err != nil {
		return err
	}
	return nil

}
