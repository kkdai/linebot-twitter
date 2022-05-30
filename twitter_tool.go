package main

import (
	"errors"

	tt "github.com/kkdai/twitter"
)

type tweet struct {
}

type Question struct {
	id     int
	text   string
	result tweet
}

func showQuestion(timeline tt.TimelineTweets) (Question, error) {
	if len(timeline) == 0 {
		return Question{}, errors.New("No Question")
	}
	return Question{}, errors.New("No Question")
}
