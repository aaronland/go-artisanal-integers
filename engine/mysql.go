package engine

// http://code.flickr.com/blog/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/
// https://github.com/go-sql-driver/mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thisisaaronland/go-artisanal-integers"
)

type MySQLEngine struct {
	artisanalinteger.Engine
	dsn string
}

func (eng *MySQLEngine) Set(i int64) error {

	max, err := eng.Max()

	if err != nil {
		return err
	}

	if i < max {
		return errors.New("integer value too small")
	}

	db, err := eng.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	sql := fmt.Sprintf("ALTER TABLE integers AUTO_INCREMENT=%d", i)

	_, err = db.Query(sql)

	/*
		st, err := db.Prepare("ALTER TABLE integers AUTO_INCREMENT=?")

		if err != nil {
			return err
		}

		defer st.Close()

		_, err = st.Exec(i)
	*/

	if err != nil {
		return err
	}

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
	rows.Next()

	var max int64

	err = rows.Scan(&max)

	if err != nil {
		return -1, err
	}

	return max, nil
}

// https://dev.mysql.com/doc/refman/5.7/en/getting-unique-id.html

func (eng *MySQLEngine) Next() (int64, error) {

	err := eng.set_autoincrement()

	if err != nil {
		return -1, err
	}

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	st_replace, err := db.Prepare("REPLACE INTO integers (stub) VALUES(?)")

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
	rows.Next()

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
		dsn: dsn,
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
