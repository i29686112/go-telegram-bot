package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	checkMigration()
	db := getDb()

	router.GET("/insert", func(c *gin.Context) {

		randomNumber := rand.Int()
		user := User{Name: "ian" + strconv.Itoa(randomNumber), Money: rand.Float32()}

		insertResult := db.Create(&user) // pass pointer of data to Create

		if insertResult.Error != nil {
			fmt.Println(insertResult.Error)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": user,
		})

	})

	router.POST("/", botMainHandler)

	err := router.Run(":80")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:80
}

func botMainHandler(c *gin.Context) {

	db := getDb()

	// it seems we can't bind after getRawData, it will cause Bind get EOF error.
	//rawData, err := c.GetRawData()
	//if err != nil {
	//	//Handle Error
	//}

	var rawJsonString = ""

	var telegramRequestBody TelegramRequestBody

	err := c.Bind(&telegramRequestBody)
	if err != nil {
		return
	}

	// insert log
	telegramWebhookHistory := saveTelegramLog(telegramRequestBody, rawJsonString, db)

	// reply user
	replyUser(telegramWebhookHistory)

	c.JSON(http.StatusOK, gin.H{
		"message": "You input is " + telegramRequestBody.Message.Text,
	})

}

func replyUser(telegramWebhookHistory TelegramWebhookHistory) {

	chatIdString := strconv.Itoa(telegramWebhookHistory.ChatId)
	replyStr := "You input string is `" + telegramWebhookHistory.MessageText + "`"
	webhookStr := "bot" + telegramBotToken

	webhookUrl := fmt.Sprintf("https://api.telegram.org/%s/sendMessage?chat_id=%s&text=%s", webhookStr, chatIdString, replyStr)

	res, err := http.Get(webhookUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	telegramSendMessageResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", telegramSendMessageResponse)
}

func saveTelegramLog(telegramRequestBody TelegramRequestBody, rawJsonString string, db *gorm.DB) TelegramWebhookHistory {

	telegramWebhookHistory := TelegramWebhookHistory{
		ChatId:      telegramRequestBody.Message.Chat.Id,
		UserId:      telegramRequestBody.Message.From.Id,
		FirstName:   telegramRequestBody.Message.From.FirstName,
		LastName:    telegramRequestBody.Message.From.LastName,
		Username:    telegramRequestBody.Message.From.Username,
		MessageDate: telegramRequestBody.Message.Date,
		MessageText: telegramRequestBody.Message.Text,
		RawRequest:  rawJsonString,
	}

	db.Create(&telegramWebhookHistory) // pass pointer of data to Create

	return telegramWebhookHistory

}

type MessageBody struct {
	Text      string   `json:"text"`
	MessageId int      `json:"message_id"`
	Date      int      `json:"date"`
	From      FromBody `json:"from"`
	Chat      FromBody `json:"chat"`
}

type FromBody struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type TelegramRequestBody struct {
	Message  MessageBody `json:"message"`
	UpdateId int         `json:"update_id"`
}
