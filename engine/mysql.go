package engine

// https://github.com/go-sql-driver/mysql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thisisaaronland/go-artisanal-integers"
	"strconv"
	"sync"
)

type MySQLEngine struct {
	artisanalinteger.Engine
	dsn string
	table string
}

func (eng *MySQLEngine) Set(i int64) error {

	db, err := &eng.Connect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

}

func (eng *MySQLEngine) Max() (int64, error) {

	db, err := &eng.Connect()

	if err != nil {
		return nil, err
	}

	defer db.Close()

}

func (eng *MySQLEngine) Next() (int64, error) {

	db, err := &eng.Connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	st, err := db.Prepare("REPLACE IN ? (stub) VALUES(?)")

	if err != nil {
		return -1, err
	}

	defer st.Close()

	_, err = st.Exec(eng.table, "a")

	if err != nil {
		return -1, err
	}

	// TO DO: get insert ID
}

func (eng *MySQLEngine) set_autoincrement() error {

	db, err := &eng.Connect()

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
