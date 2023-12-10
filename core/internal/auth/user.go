package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jyotirmoydotdev/openfy/db/models"
	database "github.com/jyotirmoydotdev/openfy/db/repositories"
	"golang.org/x/crypto/bcrypt"
)

var userIDCounter int

func RegisterUser(ctx *gin.Context) {
	var newUser struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// --------
	var newUserDatabase models.User
	err := copier.Copy(&newUserDatabase, &newUser)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// --------
	newUser.Email = strings.ToLower(newUser.Email)
	for _, u := range database.Users {
		if u.Email == newUser.Email {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already exist, please login",
			})
			return
		}
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	secretKey, err := generateRandomKey()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	}
	database.UserSecrets[newUser.Email] = secretKey
	newUserDatabase.Password = string(hashPassword)
	newUserDatabase.ID = generateUserID()
	database.Users = append(database.Users, newUserDatabase)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "User registered successfully",
		"data":   newUserDatabase,
	})
}
func LoginUser(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	loginRequest.Email = strings.ToLower(loginRequest.Email)
	var userOk bool
	if _, ok := database.UserSecrets[loginRequest.Email]; !ok {
		userOk = false
	} else {
		for _, u := range database.Users {
			if u.Email == loginRequest.Email {
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginRequest.Password)); err == nil {
					userOk = true
					break
				}
			}
		}
	}
	if userOk {
		Token, err := GenerateUserJWT(loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server error",
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"token": Token,
		})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not a valid user",
		})
		return
	}
}
func UpdateUser(ctx *gin.Context) {

}

func generateUserID() string {
	userIDCounter++
	return fmt.Sprintf("U%d", userIDCounter)
}
