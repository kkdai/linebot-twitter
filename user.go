package main

import "log"

type GameUsers struct {
	// Uid: LINE UID to storage your result and access
	// Win: How many time your guess is correct.
	// Total: Total games you played.
	// VerificationCode: Twitter ver code
	// TokenKey" Twitter tokenkey

	Uid              string `json:"uid" bson:"uid"`
	Win              int    `json:"win" bson:"win"`
	Total            int    `json:"total" bson:"total"`
	VerificationCode string `json:"vcode" bson:"vcode"`
	TokenKey         string `json:"tkey" bson:"tkey"`
}

func (u *GameUsers) Add() {
	_, err := meta.Db.Model(u).Insert()
	if err != nil {
		log.Println(err)
	}
}

func (u *GameUsers) Get(uid string) (err error) {
	log.Println("***Get user by uid=", uid)
	user := GameUsers{}
	err = meta.Db.Model(&user).
		Where("uid = ?", uid).
		Select()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("User Data DB result= ", user)

	u.Uid = user.Uid
	u.Win = user.Win
	u.Total = user.Total
	u.VerificationCode = user.VerificationCode
	u.TokenKey = user.TokenKey
	return nil
}

func (u *GameUsers) Update() (err error) {
	log.Println("***Update User=", u)

	_, err = meta.Db.Model(u).
		Set("win = ?", u.Win).
		Set("total = ?", u.Total).
		Set("vcode = ?", u.VerificationCode).
		Set("tkey = ?", u.TokenKey).
		Where("user_id = ?", u.Uid).
		Update()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
