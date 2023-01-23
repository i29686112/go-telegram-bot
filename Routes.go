package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"main/repositories"
	"main/structs"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	var telegramRequestBody structs.TelegramRequestBody

	// get user input
	body, getRawDataErr := client.GetRawData()
	if getRawDataErr != nil {
		log.Fatal(getRawDataErr)
		return
	}
	rawJsonString := string(body)

	telegramRequestBody, parseTelegramRequestBodyErr := stringParseToTelegramWebhook(rawJsonString)

	if parseTelegramRequestBodyErr != nil {
		log.Fatal(parseTelegramRequestBodyErr)
		return
	}

	// insert log
	telegramWebhookHistory := saveTelegramLog(telegramRequestBody, rawJsonString, db)

	lastUserInputCommand := repositories.GetLastUserInputCommand(telegramWebhookHistory.ChatId, db)

	if telegramWebhookHistory.MessageText == "/searchnews" {
		repositories.SaveLastTelegramCommand(telegramRequestBody, db)
	}

	// main handle function
	replyUser(telegramWebhookHistory, lastUserInputCommand)

	if lastUserInputCommand.ChatId > 0 {
		// delete last command
		repositories.DeleteLastTelegramCommand(telegramWebhookHistory.ChatId, db)
	}

	client.JSON(http.StatusOK, gin.H{
		"message": "You input is " + telegramRequestBody.Message.Text,
	})

}

func replyUser(telegramWebhookHistory TelegramWebhookHistory, lastUserInputCommand repositories.TelegramLastCommand) {

	chatIdString := strconv.Itoa(telegramWebhookHistory.ChatId)
	webhookStr := "bot" + telegramBotToken

	switch telegramWebhookHistory.MessageText {
	case "/searchnews":
		sendPlainTextToUser(webhookStr, chatIdString, "Give me a news search keyword")
	default:
		// get last command
		if lastUserInputCommand.ID > 0 {

			//
			searchResult, err := searchFromBingNewsSearch(telegramWebhookHistory.MessageText)

			if err != nil {
				return
			}

			replyText := fmt.Sprintf("No news had been found with the keyword `%s`", telegramWebhookHistory.MessageText)

			if len(searchResult.Value) == 0 {

				sendPlainTextToUser(webhookStr, chatIdString, replyText)

			} else {
				replyText = fmt.Sprintf("Here are top 3 news with the keyword `%s`", telegramWebhookHistory.MessageText)

				sendPlainTextToUser(webhookStr, chatIdString, replyText)

				for index, row := range searchResult.Value {

					if index >= 3 {
						// we only need return first 3 to user.
						break
					}
					rankText := fmt.Sprintf("======Top %d news====", index+1)
					recommendText := fmt.Sprintf("%s:\n title:%s\n link:%s", rankText, row.Name, row.Url)
					sendPlainTextToUser(webhookStr, chatIdString, recommendText)
				}

			}

		} else {
			// reply user with what they said.
			replyStr := "You input string is `" + telegramWebhookHistory.MessageText + "`"
			sendPlainTextToUser(webhookStr, chatIdString, replyStr)
		}

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
	body, err := io.ReadAll(res.Body)
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
	// parse response to struct directly
	defaultBingNewsReachResult := structs.BingNewsReachResult{}
	jsonDecodeError2 := json.NewDecoder(res.Body).Decode(&defaultBingNewsReachResult)
	if jsonDecodeError2 != nil {
		return structs.BingNewsReachResult{}, jsonDecodeError2
	}
	return defaultBingNewsReachResult, nil

}

func stringParseToTelegramWebhook(bodyToString string) (structs.TelegramRequestBody, error) {
	telegramRequestBody := structs.TelegramRequestBody{}
	err := json.Unmarshal([]byte(bodyToString), &telegramRequestBody)
	if err != nil {
		return telegramRequestBody, err
	}
	return telegramRequestBody, nil
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
	webhookUrl := fmt.Sprintf("https://api.telegram.org/%s/sendMessage", webhookStr)

	formData := url.Values{
		"chat_id": {chatIdString},
		"text":    {replyStr},
	}

	client := &http.Client{}
	req, createHttpClientErr := http.NewRequest("POST", webhookUrl, strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if createHttpClientErr != nil {
		fmt.Println(createHttpClientErr)
		return
	}

	res, requestError := client.Do(req)
	if requestError != nil {
		fmt.Println(requestError)
		return
	}
	defer res.Body.Close()
	telegramSendMessageResponse, err := io.ReadAll(res.Body)
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
