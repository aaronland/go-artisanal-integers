package service

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

type FSService struct {
	Service
	key       string
	offset    int64
	increment int64
	mu        *sync.Mutex
}

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "fs", NewFSService)

	if err != nil {
		panic(err)
	}

}

func NewFSService(ctx context.Context, uri string) (Service, error) {

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

	s := &FSService{
		key:       abs_path,
		increment: 2,
		offset:    1,
		mu:        mu,
	}

	err = SetParametersFromURI(ctx, s, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to set parameters, %w", err)
	}

	return s, nil
}

func (s *FSService) SetLastInt(ctx context.Context, i int64) error {

	last, err := s.LastInt(ctx)

	if err != nil {
		return err
	}

	if i < last {
		return errors.New("integer value too small")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return write_int(s.key, i)
}

func (s *FSService) SetOffset(ctx context.Context, i int64) error {
	s.offset = i
	return nil
}

func (s *FSService) SetIncrement(ctx context.Context, i int64) error {
	s.increment = i
	return nil
}

func (s *FSService) LastInt(ctx context.Context) (int64, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	return read_int(s.key)
}

func (s *FSService) NextInt(ctx context.Context) (int64, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	i, err := read_int(s.key)

	if err != nil {
		return -1, err
	}

	i = i + s.increment

	err = write_int(s.key, i)

	if err != nil {
		return -1, err
	}

	return i, nil
}

func (s *FSService) Close(ctx context.Context) error {
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
