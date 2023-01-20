package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	checkMigration()
	db := getDb()

	r.GET("/insert", func(c *gin.Context) {

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

	r.POST("/", func(c *gin.Context) {

		var jsonDecodeMap map[string]interface{}

		err := c.Bind(&jsonDecodeMap)
		if err != nil {
			return
		}
		//= jsonDecodeMap["message"]["text"]
		var updateId float64 = jsonDecodeMap["update_id"].(float64)
		var messageBody map[string]interface{} = jsonDecodeMap["message"].(map[string]interface{})
		var messageText string = messageBody["text"].(string)
		//var text string = jsonDecodeMap["message"]["text"].(string)

		fmt.Println(messageText)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%.0f", updateId) + ", text=" + messageText,
		})

	})

	err := r.Run(":80")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:12345
}

type TextBody struct {
	text string
}

type MessageBody struct {
	message  TextBody
	updateId int16
}
