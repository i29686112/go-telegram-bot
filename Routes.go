package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"main/structs"
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

func botMainHandler(client *gin.Context) {

	db := getDb()

	//it seems we can't bind after getRawData, it will cause Bind get EOF error.
	//rawData, err := client.GetRawData()
	//if err != nil {
	//	//Handle Error
	//}
	//fmt.Println(rawData)

	var rawJsonString = ""

	var telegramRequestBody structs.TelegramRequestBody

	err2 := client.Bind(&telegramRequestBody)
	if err2 != nil {
		return
	}

	// insert log
	telegramWebhookHistory := saveTelegramLog(telegramRequestBody, rawJsonString, db)

	replyUser(telegramWebhookHistory)

	client.JSON(http.StatusOK, gin.H{
		"message": "You input is " + telegramRequestBody.Message.Text,
	})

}

func replyUser(telegramWebhookHistory TelegramWebhookHistory) {

	chatIdString := strconv.Itoa(telegramWebhookHistory.ChatId)
	webhookStr := "bot" + telegramBotToken

	switch telegramWebhookHistory.MessageText {
	case "/searchnews":
		//
		searchResult, err := searchFromBingNewsSearch("dog")

		if err != nil {
			return
		}

		fmt.Print(searchResult)

	default:
		// reply user with what they said.
		replyStr := "You input string is `" + telegramWebhookHistory.MessageText + "`"
		sendPlainTextToUser(webhookStr, chatIdString, replyStr)

	}

}

func searchFromBingNewsSearch(searchText string) (structs.BingNewsReachResult, error) {

	defaultBingNewsReachResult := structs.BingNewsReachResult{}
	url := fmt.Sprintf("https://bing-news-search1.p.rapidapi.com/news/search?q=%s&freshness=Day&safeSearch=Off&mkt=zh-TW", searchText)
	method := "GET"

	client := &http.Client{}
	req, createHttpClientErr := http.NewRequest(method, url, nil)

	if createHttpClientErr != nil {
		fmt.Println(createHttpClientErr)
		return defaultBingNewsReachResult, nil
	}

	req.Header.Add("X-BingApis-SDK", "true")
	req.Header.Add("X-RapidAPI-Key", xRapidApiKey)
	req.Header.Add("X-RapidAPI-Host", xRapidApiHost)

	res, requestError := client.Do(req)
	if requestError != nil {
		fmt.Println(requestError)
		return defaultBingNewsReachResult, nil
	}
	defer res.Body.Close()

	// parse method 1
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return defaultBingNewsReachResult, nil
	}
	bodyToString := string(body)
	defaultBingNewsReachResult, jsonDecodeError1 := stringParseToBingSearchResult(bodyToString)
	if jsonDecodeError1 != nil {
		return structs.BingNewsReachResult{}, jsonDecodeError1
	}

	// parse method 2, we dont use it because we want gent entire response text.
	//defaultBingNewsReachResult, jsonDecodeError2 := bodyParseToBingNewsReachResult(res)
	//if jsonDecodeError2 != nil {
	//	return structs.BingNewsReachResult{}, jsonDecodeError2
	//}

	return defaultBingNewsReachResult, nil

}

func bodyParseToBingNewsReachResult(res *http.Response) (structs.BingNewsReachResult, error) {

	defaultBingNewsReachResult := structs.BingNewsReachResult{}
	jsonDecodeError2 := json.NewDecoder(res.Body).Decode(&defaultBingNewsReachResult)
	if jsonDecodeError2 != nil {
		return structs.BingNewsReachResult{}, jsonDecodeError2
	}
	return defaultBingNewsReachResult, nil

}

func stringParseToBingSearchResult(bodyToString string) (structs.BingNewsReachResult, error) {
	defaultBingNewsReachResult := structs.BingNewsReachResult{}
	err := json.Unmarshal([]byte(bodyToString), &defaultBingNewsReachResult)
	if err != nil {
		return defaultBingNewsReachResult, err
	}
	return defaultBingNewsReachResult, nil
}

func sendPlainTextToUser(webhookStr string, chatIdString string, replyStr string) {
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

func saveTelegramLog(telegramRequestBody structs.TelegramRequestBody, rawJsonString string, db *gorm.DB) TelegramWebhookHistory {

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
