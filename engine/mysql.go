package engine

// https://github.com/go-sql-driver/mysql

import (
	"database/sql"
	"database/sql/driver"
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

	st, err := db.Prepare("SELECT MAX(id) FROM ?")

	if err != nil {
		return -1, err
	}

	defer st.Close()

	row, err = st.Exec(eng.table)

	if err != nil {
		return -1, err
	}

	var max int64

	err = row.Scan(&max)

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

	st_replace, err := db.Prepare("REPLACE IN ? (stub) VALUES(?)")

	if err != nil {
		return -1, err
	}

	defer st_replace.Close()

	_, err = st_replace.Exec(eng.table, "a")

	if err != nil {
		return -1, err
	}

	st_last, err := db.Prepare("SELECT LAST_INSERT_ID()")

	if err != nil {
		return -1, err
	}

	defer st_last.Close()

	row, err = st_last.Exec()

	if err != nil {
		return -1, err
	}

	var next int64

	err = row.Scan(&next)

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

func (eng *MySQLEngine) connect() (driver.Conn, err) {

	db, err := sql.Open("mysql", eng.dsn)

	if err != nil {
		return nil, er
	}

	return db, nil
}

func NewMySQLEngine(dsn string, table string) (*MySQLEngine, error) {

	eng := MySQLEngine{
		dsn:   dsn,
		table: table,
	}

	db, err := &eng.Connect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return nil, er
	}

	return &eng, nil
}
