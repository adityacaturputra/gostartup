package main

import (
	"fmt"
	"gostartup/user"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// didapatkan dari:   https://gorm.io/docs/connecting_to_the_database
	dsn := "root:@tcp(127.0.0.1:3306)/gostartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database is good")

	var users []user.User
	length := len(users)

	fmt.Println(length)

	db.Find(&users)
	length = len(users)

	fmt.Println(length)
	for _, user := range users {
		fmt.Println("Id:", user.ID)
		fmt.Println("Name:", user.Name)
		fmt.Println("Email:", user.Email)
		fmt.Println(user.Occupation)
		fmt.Println(user.PasswordHash)
		fmt.Println(user.Role)
		fmt.Println(user.AvatarFileName)
		fmt.Println(user.CreatedAt)
		fmt.Println(user.UpdatedAt)
		fmt.Println("================")
	}
}
