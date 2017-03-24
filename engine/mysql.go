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
	_ "log"
)

type MySQLEngine struct {
	artisanalinteger.Engine
	dsn string
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

	// why doesn't this work when I use Prepare/Exec ?

	sql := fmt.Sprintf("ALTER TABLE integers AUTO_INCREMENT=%d", i)
	_, err = db.Query(sql)

	if err != nil {
		return err
	}

	return nil
}

func (eng *MySQLEngine) SetOffset(int64) error {
	return nil
}

func (eng *MySQLEngine) SetIncrement(int64) error {
	return nil
}

func (eng *MySQLEngine) LastId() (int64, error) {

	db, err := eng.connect()

	if err != nil {
		return -1, err
	}

	defer db.Close()

	row := db.QueryRow("SELECT MAX(id) FROM integers")

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

	st_replace, err := db.Prepare("REPLACE INTO integers (stub) VALUES(?)")

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
