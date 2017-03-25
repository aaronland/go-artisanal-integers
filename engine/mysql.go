package engine

// http://code.flickr.com/blog/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/
// https://github.com/go-sql-driver/mysql
// https://golang.org/pkg/database/sql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thisisaaronland/go-artisanal-integers"
)

type MySQLEngine struct {
	artisanalinteger.Engine
	dsn       string
	key       string
	offset    int64
	increment int64
}

func (eng *MySQLEngine) SetLastId(i int64) error {

	last, err := eng.LastId()

	if err != nil {
		return err
	}

	if i < last {
		return errors.New("integer value too small")
	}

	db, err := eng.connect()

	if err != nil {
		return err
	}

	defer db.Close()

	sql := fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT=%d", eng.key, i)
	st, err := db.Prepare(sql)

	if err != nil {
		return err
	}

	_, err = st.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (eng *MySQLEngine) SetKey(k string) error {
	eng.key = k
	return nil
}

func (eng *MySQLEngine) SetOffset(i int64) error {
	eng.offset = i
	return nil
}

func (eng *MySQLEngine) SetIncrement(i int64) error {
	eng.increment = i
	return nil
}

func (eng *MySQLEngine) LastId() (int64, error) {

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	sql := fmt.Sprintf("SELECT MAX(id) FROM %s", eng.key)
	row := db.QueryRow(sql)

	var max int64

	err = row.Scan(&max)

	if err != nil {
		return -1, err
	}

	return max, nil
}

// https://dev.mysql.com/doc/refman/5.7/en/getting-unique-id.html

func (eng *MySQLEngine) NextId() (int64, error) {

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	err = eng.set_autoincrement(db)

	if err != nil {
		return -1, err
	}

	sql := fmt.Sprintf("REPLACE INTO %s (stub) VALUES(?)", eng.key)
	st_replace, err := db.Prepare(sql)

	if err != nil {
		return -1, err
	}

	defer st_replace.Close()

	result, err := st_replace.Exec("a")

	if err != nil {
		return -1, err
	}

	next, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return next, nil
}

func (eng *MySQLEngine) set_autoincrement(db *sql.DB) error {

	sql_incr := fmt.Sprintf("SET @@auto_increment_increment=%d", eng.increment)
	st_incr, err := db.Prepare(sql_incr)

	if err != nil {
		return err
	}

	defer st_incr.Close()

	_, err = st_incr.Exec()

	if err != nil {
		return err
	}

	sql_off := fmt.Sprintf("SET @@auto_increment_offset=%d", eng.offset)
	st_off, err := db.Prepare(sql_off)

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
		dsn:       dsn,
		key:       "integers",
		offset:    1,
		increment: 2,
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
