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

func (g *GameData) ShowAll() (err error) {
	log.Println("***ShowAll  User -->")
	users := []GameUsers{}
	err = g.Db.Model(&users).Select()
	if err != nil {
		log.Println(err)
	}
	log.Println("***Show all users =", users)
	return nil
}

func (g *GameData) CreateSchema() error {
	models := []interface{}{
		(*GameData)(nil),
	}

	for _, model := range models {
		err := g.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}
