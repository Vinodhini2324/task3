package database

import (
	"jwtapi/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnection() {
	con, err := gorm.Open(mysql.Open("root:Yes@/jwtapi"), &gorm.Config{})
	if err != nil {
		panic("can't able to connect with database")
	}
	con.AutoMigrate(&models.Users{})
	DB = con
}
