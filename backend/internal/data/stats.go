package data

import (
	"database/sql"
	"log/slog"
)

type StatsModel struct {
	DB *sql.DB
	logger *slog.Logger
}

type WeightWeeklyAverage struct {
    Year int `json:"year"`
    Week int `json:"week"`
    Weight float64 `json:"weight"`
}

func (m *StatsModel) GetWeightDifference(userID int64, filters Filters) (*float64, error) {
    
    m.logger.Debug("Getting weight difference...")

    format := "2006-01-02"
    var args []any
    q := ""

    // If no start or end time is provided, get the difference between the first and last checkin
    if filters.StartTime.IsZero() && filters.EndTime.IsZero() {
        q = `
        SELECT
        ROUND(
            (SELECT weight FROM checkins WHERE user_id = ? ORDER BY datetime DESC LIMIT 1)
            -
            (SELECT weight FROM checkins WHERE user_id = ? ORDER BY datetime ASC LIMIT 1),
            1
        );
        `
        args = []any{userID, userID}

    // If only an end time is provided, get the difference between the first checkin and the checkin at the end time
    } else if filters.StartTime.IsZero() && !filters.EndTime.IsZero() {
        q = `
        SELECT
        ROUND(
            (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1)
            -
            (SELECT weight FROM checkins WHERE user_id = ? ORDER BY datetime ASC LIMIT 1),
            1
        );
        `
        args = []any{
            filters.EndTime.Format(format),
            filters.EndTime.AddDate(0, 0, 1).Format(format),
            userID,
            userID,
        }
    
    // If only a start time is provided, get the difference between the last checkin and the checkin at the start time
    } else if !filters.StartTime.IsZero() && filters.EndTime.IsZero() {
        q = `
        SELECT
        ROUND(
            (SELECT weight FROM checkins WHERE user_id = ? ORDER BY datetime DESC LIMIT 1)
            -
            (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime DESC LIMIT 1),
            1
        );
        `
        args = []any{
            userID,
            filters.StartTime.Format(format),
            filters.StartTime.AddDate(0, 0, 1).Format(format),
            userID,
        }

    // If a start time and end time is provided, get the weight difference between those dates.
    } else {
        q = `
        SELECT
        ROUND(
            (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1)
            -
            (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1),
            1
        );
        `
        args = []any{
            filters.EndTime.Format(format),
            filters.EndTime.AddDate(0, 0, 1).Format(format),
            userID,
            filters.StartTime.Format(format),
            filters.StartTime.AddDate(0, 0, 1).Format(format),
            userID,
        }

    }

    var weight *float64
    err := m.DB.QueryRow(q, args...).Scan(&weight)
    if err != nil {
        return nil, err
    }

    return weight, nil
}

func (m *StatsModel) GetWeightAverage(userID int64, filters Filters) (*float64, error) {
    
    m.logger.Debug("Getting weight average...")

    format := "2006-01-02"
    var args []any
    q := ""

    if filters.StartTime.IsZero() {
        q = `
        SELECT ROUND(AVG(weight), 1) FROM checkins WHERE user_id = ? AND datetime < strftime('%s', ?);
        `
        args = []any{
            userID,
            filters.EndTime.AddDate(0,0,1).Format(format),
        }
    } else {
        q = `
        SELECT ROUND(AVG(weight), 1) FROM checkins WHERE user_id = ? AND datetime >= strftime('%s', ?) AND datetime < strftime('%s', ?);
        `
        args = []any{
            userID,
            filters.StartTime.Format(format),
            filters.EndTime.AddDate(0,0,1).Format(format),
        }
    }

    var weight *float64
    err := m.DB.QueryRow(q, args...).Scan(&weight)
    if err != nil {
        return nil, err
    }

    return weight, nil
}

func (m *StatsModel) GetWeightAverageByWeek(userID int64) ([]WeightWeeklyAverage, error) {
        
        m.logger.Debug("Getting weight average by week...")

        q := `
        SELECT strftime('%Y', datetime(datetime, 'unixepoch')) as year, strftime('%W', datetime(datetime, 'unixepoch')) as week, ROUND(AVG(weight), 1) 
        FROM checkins 
        WHERE user_id = ?
        GROUP BY strftime('%Y', datetime(datetime, 'unixepoch')), strftime('%W', datetime(datetime, 'unixepoch'))
        ORDER BY strftime('%Y', datetime(datetime, 'unixepoch')), strftime('%W', datetime(datetime, 'unixepoch'));
        `
        args := []any{
            userID,
        }
    
        // Execute query
        rows, err := m.DB.Query(q, args...)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        // Parse data
        data := []WeightWeeklyAverage{}
        for rows.Next() {
            var r WeightWeeklyAverage
            err := rows.Scan(&r.Year, &r.Week, &r.Weight)
            if err != nil {
                return nil, err
            }
            data = append(data, r)
        }
    
        return data, nil
    }
