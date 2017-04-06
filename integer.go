package artisanalinteger

type Engine interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
	SetLastInt(int64) error
	SetKey(string) error
	SetOffset(int64) error
	SetIncrement(int64) error
}

type Service interface {
	NextInt() (int64, error)
	LastInt() (int64, error)
}

// EXPERIMENTAL

type Server interface {
	ListenAndServer(Service) error
}
