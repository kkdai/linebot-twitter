package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// GetTwitterToken:
func GetTwitterToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enter Get twitter token")
	values := r.URL.Query()
	user.VerificationCode = values.Get("oauth_verifier")
	user.TokenKey = values.Get("oauth_token")

	if err := user.Update(); err != nil {
		log.Println("Update user failed, ", err)
	}

	GetQuestion()
}

// RedirectUserToTwitter
func GetTwitterURL() string {
	requestUrl, _ := twitterClient.GetAuthURL(CallbackURL)
	log.Println("CallbackURL=", CallbackURL, " requestUrl url=", requestUrl)
	return requestUrl
}

func GetQuestion() {
	if len(user.VerificationCode) == 0 || len(user.TokenKey) == 0 {
		return
	}

	// Complete twitter auth.
	twitterClient.CompleteAuth(user.TokenKey, user.VerificationCode)

	// Get timeline
	timeline, _, _ := twitterClient.QueryTimeLine(1)
	ret := fmt.Sprintf("TimeLine \n\n=%v", timeline)
	// fmt.Fprintf(w, ret+" \n\n The item is: "+string(bits))

	// push message
	if _, err := bot.PushMessage(user.Uid, linebot.NewTextMessage("Timeline \n"+ret)).Do(); err != nil {
		log.Print(err)
	}
}
