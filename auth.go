package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

//temp
var verificationCode string
var tokenKey string

// GetTwitterToken:
func GetTwitterToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter Get twitter token")
	values := r.URL.Query()
	verificationCode = values.Get("oauth_verifier")
	tokenKey = values.Get("oauth_token")

	SendQuestion()
}

// RedirectUserToTwitter
func GetTwitterURL() string {
	requestUrl, _ := twitterClient.GetAuthURL(CallbackURL)
	log.Println("CallbackURL=", CallbackURL, " requestUrl url=", requestUrl)
	return requestUrl
}

func SendQuestion() {
	if len(verificationCode) == 0 || len(verificationCode) == 0 {
		return
	}

	// Complete twitter auth.
	twitterClient.CompleteAuth(tokenKey, verificationCode)

	// Get timeline
	timeline, _, _ := twitterClient.QueryTimeLine(1)
	ret := fmt.Sprintf("TimeLine \n\n=%v", timeline)
	// fmt.Fprintf(w, ret+" \n\n The item is: "+string(bits))

	// push message
	if _, err := bot.PushMessage(user.Uid, linebot.NewTextMessage("Timeline \n"+ret)).Do(); err != nil {
		log.Print(err)
	}
}
