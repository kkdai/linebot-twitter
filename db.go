package main

import (
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type GameData struct {
	Db  *pg.DB
	Log *log.Logger
}

type GameUser struct {
	uid   string
	win   int
	total int
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*GameUser)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
