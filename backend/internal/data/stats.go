package data

import (
	"database/sql"
	// "errors"
	"log/slog"
)

type StatsModel struct {
	DB *sql.DB
	logger *slog.Logger
}

type Stats struct {
	WeightDifference WeightDifferenceStats 	`json:"weight_difference"`
}

type WeightDifferenceStats struct {
	AllTime float64 		`json:"all_time"`
	WeekAgo float64 		`json:"week_ago"`
	NinetyDaysAgo float64 	`json:"ninety_days_ago"`
}

func (m *StatsModel) GetWeightDifference(userID int64, filters Filters) (*float64, error) {
    
    m.logger.Debug("Getting weight difference...")


    format := "2006-01-02"
    args := []any{}
    q := ""

    if filters.StartTime.IsZero() {
        q = `
        SELECT
        (SELECT weight FROM checkins WHERE user_id = ? ORDER BY datetime ASC LIMIT 1)
        -
        (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1);
        `
        args = []any{
            userID,
            filters.EndTime.Format(format),
            filters.EndTime.AddDate(0, 0, 1).Format(format),
            userID,
        }
    } else {
        q = `
        SELECT
        (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1)
        -
        (SELECT weight FROM checkins WHERE datetime > strftime('%s', ?) AND datetime < strftime('%s', ?) AND user_id = ? ORDER BY datetime LIMIT 1);
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

// func (m *StatsModel) GetStats(userID int64) (*Stats, error) {

// 	m.logger.Debug("Getting stats...")

//     // TODO: Simplify this query.
// 	q := `
// 	SELECT 
//     wdat1.weight - wdat2.weight AS weight_difference_all_time, 
//     wdsd1.weight - wdsd2.weight AS weight_difference_7_days,
//     CASE WHEN (
//             SELECT datetime('now') - MIN(datetime)
//             FROM checkins
//             WHERE datetime >= strftime('%s', datetime('now','-91 day')) AND user_id=?
//         ) < 91 THEN 0
//         ELSE (
//             SELECT wdnd1.weight - wdnd2.weight
//             FROM
//                 ( SELECT weight
//                     FROM checkins
//                     WHERE datetime = (
//                              SELECT
//                               MAX(datetime)
//                                 FROM checkins
//                              WHERE datetime >= strftime('%s', datetime('now','-91 day')) AND user_id=?
//                         )
//                 ) AS wdnd1,
//                 (
//                     SELECT
//                         weight
//                     FROM
//                         checkins
//                     WHERE
//                         datetime = (
//                             SELECT
//                                 MIN(datetime)
//                             FROM
//                                 checkins
//                             WHERE
//                                 datetime >= strftime('%s', datetime('now','-91 day')) AND user_id=?
//                         )
//                 ) AS wdnd2
//         )
//     END AS weight_difference_90_days
// FROM
//     (
//         SELECT weight 
//         FROM checkins 
//         WHERE datetime = (
//             SELECT MAX(datetime) 
//             FROM checkins WHERE user_id=?
//         ) AND user_id=?
//     ) AS wdat1,
//     (
//         SELECT weight 
//         FROM checkins 
//         WHERE datetime = (
//             SELECT MIN(datetime) 
//             FROM checkins WHERE user_id=?
//         ) AND user_id=?
//     ) AS wdat2,
//     (
//         SELECT weight 
//         FROM checkins 
//         WHERE datetime = (
//             SELECT MAX(datetime) 
//             FROM checkins 
//             WHERE datetime >= strftime('%s', datetime('now','-8 day')) AND user_id=?
//         ) AND user_id=?
//     ) AS wdsd1,
//     (
//         SELECT weight 
//         FROM checkins 
//         WHERE datetime = (
//             SELECT MIN(datetime) 
//             FROM checkins 
//             WHERE datetime >= strftime('%s', datetime('now','-8 day')) AND user_id=?
//         ) AND user_id=?
//     ) AS wdsd2;
// 	`

// 	var stats Stats
// 	err := m.DB.QueryRow(q, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID).Scan(&stats.WeightDifference.AllTime, &stats.WeightDifference.WeekAgo, &stats.WeightDifference.NinetyDaysAgo)
    
// 	if err != nil {
// 		switch {
//             case errors.Is(err, sql.ErrNoRows):
//                 return &stats, nil
//             default:
//                 return nil, err
// 		}
// 	}

// 	return &stats, nil
// }