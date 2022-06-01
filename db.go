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

type GameUsers struct {
	// LINE UID to storage your result and access
	uid string

	// How many time your guess is correct.
	win int

	// Total games you played.
	total int
}

func (u *GameUsers) Add(meta *GameData) {
	_, err := meta.Db.Model(u).Insert()
	if err != nil {
		meta.Log.Println(err)
	}
}

func (u *GameUsers) Get(meta *GameData) (result *GameUsers, err error) {
	log.Println("***Get Fav uUID=", u.uid)
	userFav := GameUsers{}
	err = meta.Db.Model(&userFav).
		Where("user_id = ?", u.uid).
		Select()
	if err != nil {
		meta.Log.Println(err)
		return nil, err
	}
	meta.Log.Println("UserFavorite DB result= ", userFav)
	return &userFav, nil
}

func (u *GameUsers) Update(meta *GameData) (err error) {
	log.Println("***Update Fav User=", u)

	_, err = meta.Db.Model(u).
		Set("favorites = ?", u.uid).
		Where("user_id = ?", u.win).
		Update()
	if err != nil {
		meta.Log.Println(err)
	}
	return nil
}

func (u *GameUsers) ShowAll(meta *GameData) (err error) {
	log.Println("***ShowAll  User -->")
	users := []GameUsers{}
	err = meta.Db.Model(&users).Select()
	if err != nil {
		log.Println(err)
	}
	log.Println("***Show all users =", users)
	return nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*GameData)(nil),
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
