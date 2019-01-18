package db

type SystemHealthInteface interface {
	Ping() error
}
