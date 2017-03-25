package artisanalinteger

type Engine interface {
	NextId() (int64, error)
	LastId() (int64, error)
	SetLastId(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
}

type Service interface {
	NextId() (int64, error)
	MaxId() (int64, error)
}
