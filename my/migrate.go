package my

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Migrate program.
func Migrate() {
	dsn := "root:81202@Musuke@tcp(127.0.0.1:3306)/data_splite3?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	db.AutoMigrate(&User{}, &Group{}, &Post{}, &Comment{})
}
