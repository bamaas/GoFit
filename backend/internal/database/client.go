package database

import (
	"log/slog"

	"database/sql"
	_ "modernc.org/sqlite"
)

type CheckIn struct {
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
}

type Database struct {
	*sql.DB
	logger *slog.Logger
}

func New(logger *slog.Logger) (*Database, error){

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
	id INTEGER NOT NULL PRIMARY KEY,
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

func parseRowsToEntries(r *sql.Rows) ([]CheckIn, error){
	entries := []CheckIn{}
	for r.Next() {
		var e CheckIn
		err := r.Scan(&e.ID, &e.Weight)
		if err != nil {
			return []CheckIn{}, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (d *Database) GetCheckIn(id int) (CheckIn, error) {

	d.logger.Debug("Get entry", "id", id)

	q := `
	SELECT id, weight
	FROM entries
	WHERE id=?`

	r, err := d.Query(q, id)
	if err != nil {
		return CheckIn{}, err
	}

	entries, err := parseRowsToEntries(r)
	if err != nil {
		return CheckIn{}, err
	}

	return entries[0], nil
}

func (d *Database) GetCheckIns() ([]CheckIn, error) {

	d.logger.Debug("Get all the entries")

	q := `
	SELECT id, weight
	FROM entries
	`
	r, err := d.Query(q)
	if err != nil {
		return nil, err
	}

	return parseRowsToEntries(r)
}

func (d *Database) InsertCheckIn(e CheckIn) error {

	d.logger.Debug("Insert check-in", "check-in", e)

	q := `
	INSERT INTO entries
	(id, weight)
	VALUES
	(?, ?);
	`
	_, err := d.Exec(q, e.ID, e.Weight)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteCheckIn(id int) error {

	d.logger.Debug("Deleting entry", "id", id)

	q := `
	DELETE FROM entries
	WHERE
	id=?
	`
	_, err := d.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateCheckIn(e CheckIn) error {

	d.logger.Debug("Updating check-in", "id", e.ID)

	q := `
	UPDATE entries
	SET weight=?
	WHERE id=?
	`
	_, err := d.Exec(q, e.Weight, e.ID)
	if err != nil {
		return err
	}
	return nil

}
