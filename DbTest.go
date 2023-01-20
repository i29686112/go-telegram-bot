package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name         string  `gorm:"size:255;index:idx_name"`
	Money        float32 `gorm:"index:idx_money"`
	NewFieldTest string
}

func getDb() *gorm.DB {
	dsn := "host=localhost user=ian password=secret dbname=bot port=5432 sslmode=disable TimeZone=Asia/Taipei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func checkMigration() {

	db := getDb()
	err := db.AutoMigrate(&User{})
	err = db.Migrator().DropIndex(&User{}, "idx_name")
	if err != nil {
		log.Fatalln("migration with error:" + err.Error())
	}
}
