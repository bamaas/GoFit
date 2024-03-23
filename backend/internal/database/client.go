package database

import (
	"log/slog"

	"database/sql"
	_ "modernc.org/sqlite"
)

type Entry struct {
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
}

type db struct {
	*sql.DB
	logger *slog.Logger
}

func New(logger *slog.Logger) (*db, error){

	logger.Debug("Intializing database...")

	d, err := sql.Open("sqlite", "gofit.db")
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

	return &db{
		d,
		logger,
	}, nil
}

func parseRowsToEntries(r *sql.Rows) ([]Entry, error){
	entries := []Entry{}
	for r.Next() {
		var e Entry
		err := r.Scan(&e.ID, &e.Weight)
		if err != nil {
			return []Entry{}, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (d *db) GetEntry(id int) (Entry, error) {

	d.logger.Debug("Get entry", "id", id)

	q := `
	SELECT id, weight 
	FROM entries 
	WHERE id=?`

	r, err := d.Query(q, id)
	if err != nil {
		return Entry{}, err
	}

	entries, err := parseRowsToEntries(r)
	if err != nil {
		return Entry{}, err
	}

	return entries[0], nil
}

func (d *db) GetEntries() ([]Entry, error) {

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

func (d *db) InsertEntry(e Entry) error {

	d.logger.Debug("Insert entry", "entry", e)

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

func (d *db) DeleteEntry(id int) error {

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

func (d *db) UpdateEntry(e Entry) error {

	d.logger.Debug("Updating entry", "id", e.ID)

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