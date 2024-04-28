package data

import (
	"database/sql"
	"errors"
	"log/slog"
)

var (
	ErrRecordNotFound = errors.New("record not found");
)

type Models struct {
	CheckIns 	CheckInModel
	Users		UserModel
	Stats		StatsModel
	Tokens		TokenModel
}

func NewModels(db *sql.DB, logger *slog.Logger) Models {
	return Models{
		CheckIns: CheckInModel{db, logger},
		Users: UserModel{db, logger},
		Stats: StatsModel{db, logger},
		Tokens: TokenModel{db, logger},
	}
}