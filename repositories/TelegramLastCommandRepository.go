package repositories

import (
	"gorm.io/gorm"
	"main/structs"
)

func GetLastUserInputCommand(chatId int, db *gorm.DB) TelegramLastCommand {
	telegramLastCommand := TelegramLastCommand{}
	db.First(&telegramLastCommand, "chat_id = ?", chatId)
	return telegramLastCommand
}

func SaveLastTelegramCommand(telegramRequestBody structs.TelegramRequestBody, db *gorm.DB) TelegramLastCommand {

	telegramLastCommand := TelegramLastCommand{
		ChatId:      telegramRequestBody.Message.Chat.Id,
		UserId:      telegramRequestBody.Message.From.Id,
		LastCommand: telegramRequestBody.Message.Text,
	}

	db.Create(&telegramLastCommand) // pass pointer of data to Create

	return telegramLastCommand

}

func DeleteLastTelegramCommand(chatId int, db *gorm.DB) TelegramLastCommand {

	telegramLastCommand := TelegramLastCommand{}

	db.Delete(&telegramLastCommand, "chat_id = ?", chatId)

	return telegramLastCommand

}

type TelegramLastCommand struct {
	// force makes the ID as the first column.
	ID          uint `gorm:"primarykey"`
	ChatId      int  `gorm:"index:idx_telegram_chat_id"`
	UserId      int  `gorm:"index:idx_telegram_user_id"`
	LastCommand string
	gorm.Model
}
