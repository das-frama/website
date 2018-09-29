package app

import (
	"log"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var Session *r.Session

// OpenSession открывает и возвращает указатель на объект сессии RethinkDB на базе полученного конфига.
func OpenSession(config *Config) (*r.Session, error) {
	log.Println("Database is connecting...")
	var err error
	Session, err = r.Connect(r.ConnectOpts{
		Address:  config.DbAddress,
		Database: config.DbName,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	log.Println("Database is connected.")

	return Session, nil
}
