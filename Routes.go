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

		var messageBody MessageBody

		err := c.Bind(&messageBody)
		if err != nil {
			return
		}

		fmt.Println(messageBody.Message.Text)

		c.JSON(http.StatusOK, gin.H{
			"message": "You input is " + messageBody.Message.Text,
		})

	})

	err := r.Run(":80")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:12345
}

type TextBody struct {
	Text string `json:"text"`
}

type MessageBody struct {
	Message  TextBody `json:"message"`
	UpdateId int      `json:"update_id"`
}
