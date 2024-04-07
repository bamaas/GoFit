package database

import (
	"log/slog"
	"time"

	"database/sql"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type CheckIn struct {
	UUID   		string    `json:"uuid,omitempty"`
	Datetime   	time.Time `json:"datetime"`
	Weight 		float64   `json:"weight"`
}

type Database struct {
	*sql.DB
	logger *slog.Logger
}

func New(logger *slog.Logger) (*Database, error) {

	logger.Debug("Intializing database...")

	d, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	err = d.Ping()
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS entries (
	uuid STRING NOT NULL PRIMARY KEY,
	datetime INTEGER NOT NULL,
	weight FLOAT NOT NULL
	);`
	_, err = d.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	// Insert some data
	return &Database{
		d,
		logger,
	}, nil
}

func (d *Database) InjectSampleData() error {

	var checkIns []CheckIn

	for i := 1; i <= 30; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		checkIn := CheckIn{
			UUID:   	uuid.String(),
			Datetime:   time.Now().AddDate(0, 0, -i),
			Weight: 	float64(i + 29),
		}
		checkIns = append(checkIns, checkIn)
	}

	for _, c := range checkIns {
		err := d.InsertCheckIn(c)
		if err != nil {
			return err
		}
	}
	return nil

}

func parseRowsToEntries(r *sql.Rows) ([]CheckIn, error) {

	// Parse db sql rows
	type dbRow struct {
		UUID   string
		Datetime   int64
		Weight float64
	}
	dbData := []dbRow{};
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
			UUID:   dbData[i].UUID,
			Datetime:   time.Unix(dbData[i].Datetime, 0),
			Weight: dbData[i].Weight})
	}

	return entries, nil
}

func (d *Database) GetCheckIn(UUID string) (CheckIn, error) {

	d.logger.Debug("Get entry", "UUID", UUID)

	q := `
	SELECT uuid, datetime, weight
	FROM entries
	WHERE uuid=?`

	r, err := d.Query(q, UUID)
	if err != nil {
		return CheckIn{}, err
	}

	entries, err := parseRowsToEntries(r)
	if err != nil {
		return CheckIn{}, err
	}

	return entries[0], nil
}

func (d *Database) GetCheckIns(filters Filters) ([]CheckIn, error) {

	d.logger.Debug("Get all the entries")

	q := `
	SELECT uuid, datetime, weight
	FROM entries
	ORDER BY datetime DESC
	LIMIT ?
	OFFSET ?;
	`
	r, err := d.Query(q, filters.limit(), filters.offset())
	if err != nil {
		return nil, err
	}

	return parseRowsToEntries(r)
}

func (d *Database) InsertCheckIn(checkIn CheckIn) error {

	d.logger.Debug("Insert check-in", "check-in", checkIn)

	q := `
	INSERT INTO entries
	(uuid, datetime, weight)
	VALUES
	(?, ?, ?);
	`
	_, err := d.Exec(q, checkIn.UUID, checkIn.Datetime.Unix(), checkIn.Weight)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteCheckIn(UUID string) error {

	d.logger.Debug("Deleting", "UUID", UUID)

	q := `
	DELETE FROM entries
	WHERE
	uuid=?
	`
	_, err := d.Exec(q, UUID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateCheckIn(checkIn CheckIn) error {

	d.logger.Debug("Updating", "UUID", checkIn.UUID)

	q := `
	UPDATE entries
	SET weight=?, datetime=?
	WHERE uuid=?
	`
	_, err := d.Exec(q, checkIn.Weight, checkIn.Datetime.Unix(), checkIn.UUID)
	if err != nil {
		return err
	}
	return nil

}
