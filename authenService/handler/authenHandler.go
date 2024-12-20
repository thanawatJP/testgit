package handler

import (
	"authenservice/database"
	"authenservice/database/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var googleOauthConfig = &oauth2.Config{}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleStartHandler(c *gin.Context) {
	// สร้าง URL สำหรับ Google OAuth
	authURL := googleOauthConfig.AuthCodeURL("random-state", oauth2.AccessTypeOffline)

	// Redirect ผู้ใช้ไปยัง Google OAuth URL
	c.Redirect(http.StatusFound, authURL)
}

func GoogleCallbackHandler(c *gin.Context) {
	var user models.UserAuth
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		fmt.Println("Error exchanging code: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println("Error getting user info: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := userInfo["email"].(string)
	result := database.DB.Where("Email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			user = models.UserAuth{
				Email:     email,
				FirstName: userInfo["given_name"].(string),
				LastName:  userInfo["family_name"].(string),
				Password:  "1234",
				RoleID:    2,
			}
			if _, err := RegisterUser(user); err != nil {
				statusCode := http.StatusInternalServerError
				if errors.Is(err, gorm.ErrRecordNotFound) {
					statusCode = http.StatusBadRequest
				}
				c.JSON(statusCode, gin.H{"error": err.Error()})
				return
			}
		} else {
			fmt.Println("Error querying user: " + result.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	signedToken, err := SignJWT(userInfo)
	if err != nil {
		fmt.Println("Error signing token: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken, "userInfo": userInfo})
}

// ยังไม่ได้ยืนยัน password ว่าถูกต้องหรือไม่
func NormalAuthHandler(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.UserAuth
	result := database.DB.Where("Email = ?", loginRequest.Email).Preload("Role").First(&user)
	if result.Error != nil {
		fmt.Println("Error querying user: " + result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		// ถ้ารหัสผ่านไม่ตรงกัน
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	userInfo := map[string]interface{}{
		"id":    user.ID,
		"name":  user.FirstName + " " + user.LastName,
		"email": user.Email,
	}

	signedToken, err := SignJWT(userInfo)
	if err != nil {
		fmt.Println("Error signing token: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken, "userInfo": user})
}
