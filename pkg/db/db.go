package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// create table if not exists
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

var database *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	if err != nil && !install {
		return err
	}

	database, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err = database.Exec(schema); err != nil {
			_ = database.Close()
			database = nil
			return err
		}
	}

	return nil
}

func Close() error {
	if database == nil {
		return nil
	}

	err := database.Close()
	database = nil
	return err
}
