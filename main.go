// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	tt "github.com/kkdai/twitter"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client
var ConsumerKey string
var ConsumerSecret string
var CallbackURL string
var twitterClient *tt.ServerClient
var meta = &GameData{}
var user = &GameUsers{}

func init() {
	//Twitter Dev Info from https://developer.twitter.com/en/apps
	ConsumerKey = os.Getenv("ConsumerKey")
	ConsumerSecret = os.Getenv("ConsumerSecret")

	//This URL need note as follow:
	// 1. Could not be localhost, change your hosts to a specific domain name
	// 2. This setting must be identical with your app setting on twitter Dev
	// 3. It should be present as "http://YOURDOMAIN.com/maketoken"
	CallbackURL = os.Getenv("CallbackURL")
}

func main() {
	var err error

	// Init LINEBot client
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)

	// Init twttr client
	fmt.Println("[app] Init server key=", ConsumerKey, " secret=", ConsumerSecret)
	twitterClient = tt.NewServerClient(ConsumerKey, ConsumerSecret)

	// API entry
	http.HandleFunc("/maketoken", GetTwitterToken)
	http.HandleFunc("/callback", callbackHandler)

	// DB Init
	dbURL := os.Getenv("DATABASE_URL")
	options, _ := pg.ParseURL(dbURL)
	db := pg.Connect(options)
	meta.Db = db
	defer db.Close()

	// Create DB if not exist.
	if err = meta.CreateSchema(); err != nil {
		panic(err)
	}

	// List all user when start.
	if err = meta.ShowAll(); err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}

				if message.Text == "auth" {
					user.Uid = event.Source.UserID
					log.Println("UID =", user.Uid)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("準備認證: "+GetTwitterURL())).Do(); err != nil {
						log.Print(err)
					}
				} else {
					// message.ID: Msg unique ID
					// message.Text: Msg text
					SendQuestion()

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
						log.Print(err)
					}
				}

			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
