package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"main/repositories"
)

type User struct {
	gorm.Model
	Name         string  `gorm:"size:255;index:idx_name"`
	Money        float32 `gorm:"index:idx_money"`
	NewFieldTest string
}

type TelegramWebhookHistory struct {
	// force makes the ID as the first column.
	ID          uint `gorm:"primarykey"`
	ChatId      int  `gorm:"index:idx_telegram_chat_id"`
	UserId      int  `gorm:"index:idx_telegram_user_id"`
	FirstName   string
	LastName    string
	Username    string
	MessageDate int
	MessageText string
	RawRequest  string
	Latitude    string `gorm:"size:20"`
	Longitude   string `gorm:"size:20"`

	gorm.Model
}

func getDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Taipei", dbHost, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func checkMigration() {

	db := getDb()
	err := db.AutoMigrate(&User{}, &TelegramWebhookHistory{}, &repositories.TelegramLastCommand{})
	//err = db.Migrator().DropIndex(&User{}, "idx_name")
	if err != nil {
		log.Fatalln("migration with error:" + err.Error())
	}
}
