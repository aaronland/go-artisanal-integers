package database

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type FSDatabase struct {
	Database
	key       string
	offset    int64
	increment int64
	mu        *sync.Mutex
}

func init() {

	ctx := context.Background()
	err := RegisterDatabase(ctx, "fs", NewFSDatabase)

	if err != nil {
		panic(err)
	}

}

func NewFSDatabase(ctx context.Context, uri string) (Database, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	abs_path, err := filepath.Abs(u.Path)

	if err != nil {
		return nil, err
	}

	root := filepath.Dir(abs_path)

	_, err = os.Stat(root)

	if os.IsNotExist(err) {

		err := os.MkdirAll(root, 0755)

		if err != nil {
			return nil, err
		}
	}

	_, err = os.Stat(abs_path)

	if os.IsNotExist(err) {

		err := write_int(abs_path, 0)

		if err != nil {
			return nil, err
		}
	}

	mu := new(sync.Mutex)

	db := &FSDatabase{
		key:       abs_path,
		increment: 2,
		offset:    1,
		mu:        mu,
	}

	err = SetParametersFromURI(ctx, db, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to set parameters, %w", err)
	}

	return db, nil
}

func (db *FSDatabase) SetLastInt(ctx context.Context, i int64) error {

	last, err := db.LastInt(ctx)

	if err != nil {
		return err
	}

	if i < last {
		return errors.New("integer value too small")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	return write_int(db.key, i)
}

func (db *FSDatabase) SetOffset(ctx context.Context, i int64) error {
	db.offset = i
	return nil
}

func (db *FSDatabase) SetIncrement(ctx context.Context, i int64) error {
	db.increment = i
	return nil
}

func (db *FSDatabase) LastInt(ctx context.Context) (int64, error) {

	db.mu.Lock()
	defer db.mu.Unlock()

	return read_int(db.key)
}

func (db *FSDatabase) NextInt(ctx context.Context) (int64, error) {

	db.mu.Lock()
	defer db.mu.Unlock()

	i, err := read_int(db.key)

	if err != nil {
		return -1, err
	}

	i = i + db.increment

	err = write_int(db.key, i)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func (db *FSDatabase) Close(ctx context.Context) error {
	return nil
}

func read_int(path string) (int64, error) {

	fh, err := os.Open(path)

	if err != nil {
		return -1, err
	}

	defer fh.Close()

	b, err := io.ReadAll(fh)

	if err != nil {
		return -1, err
	}

	i, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func write_int(path string, i int64) error {

	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer fh.Close()

	body := fmt.Sprintf("%d", i)

	_, err = fh.Write([]byte(body))
	return err
}
