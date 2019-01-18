package model

type SystemHealthInteface interface {
	PingDatabase() error
}

func (model *Model) PingDatabase() error {
	return model.db.Ping()
}
