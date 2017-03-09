package artisanalinteger

type Engine interface {
	Next() (int64, error)
	Max() (int64, error)
	Set(int64) error
}

type Integer interface {
	Next() (int64, error)
	Max() (int64, error)
}
