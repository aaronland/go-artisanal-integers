package engine

// https://github.com/go-sql-driver/mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thisisaaronland/go-artisanal-integers"
)

type MySQLEngine struct {
	artisanalinteger.Engine
	dsn   string
	table string
}

func (eng *MySQLEngine) Set(i int64) error {

	db, err := eng.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	// FIX ME
	return nil
}

func (eng *MySQLEngine) Max() (int64, error) {

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT MAX(id) FROM integers")

	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var max int64

	err = rows.Scan(&max)

	if err != nil {
		return -1, err
	}

	return max, nil
}

// https://dev.mysql.com/doc/refman/5.7/en/getting-unique-id.html

func (eng *MySQLEngine) Next() (int64, error) {

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	st_replace, err := db.Prepare("REPLACE IN integers (stub) VALUES(?)")

	if err != nil {
		return -1, err
	}

	defer st_replace.Close()

	_, err = st_replace.Exec("a")

	if err != nil {
		return -1, err
	}

	rows, err := db.Query("SELECT LAST_INSERT_ID()")

	defer rows.Close()

	var next int64

	err = rows.Scan(&next)

	if err != nil {
		return -1, err
	}

	return next, nil
}

func (eng *MySQLEngine) set_autoincrement() error {

	db, err := eng.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	st_incr, err := db.Prepare("SET @@auto_increment_increment=2")

	if err != nil {
		return err
	}

	defer st_incr.Close()

	_, err = st_incr.Exec()

	if err != nil {
		return err
	}

	st_off, err := db.Prepare("SET @@auto_increment_offset=1")

	if err != nil {
		return err
	}

	defer st_off.Close()

	_, err = st_off.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (eng *MySQLEngine) connect() (*sql.DB, error) {

	db, err := sql.Open("mysql", eng.dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMySQLEngine(dsn string) (*MySQLEngine, error) {

	eng := &MySQLEngine{
		dsn:   dsn,
		table: "integers",
	}

	db, err := eng.connect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return eng, nil
}
