package main

import (
	"fmt"
	"gostartup/auth"
	"gostartup/campaign"
	"gostartup/handler"
	"gostartup/helper"
	"gostartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/gostartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)

	campaignRepository := campaign.NewRepository(db)
	campaigns, err := campaignRepository.FindAll()

	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("jumlah seluruh campaign", len(campaigns))
	if err != nil {
		fmt.Println(err, campaigns)
	}
	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
	}

	userCampaigns, err := campaignRepository.FindByUserID(2)

	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	if err != nil {
		fmt.Println(err, userCampaigns)
	}
	fmt.Println("jumlah seluruh campaign user 2", len(userCampaigns))
	for _, campaign := range userCampaigns {
		fmt.Println(campaign.Name)
	}
	return

	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

// ambil nilai header Authorization
// dari header auhorization, kita ambil tokennya saja
// kita validasi token
// kita ambil user_id
// ambil user dari db berdasarkan user_id lewat service
// kita set context isinya user
