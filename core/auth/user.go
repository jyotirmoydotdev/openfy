package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var users []User

var userSecrets = make(map[string]string)

func RegisterUser(ctx *gin.Context) {
	var newUser User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	for _, u := range users {
		if u.Username == newUser.Username {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "User name not available",
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
	userSecrets[newUser.Username] = secretKey
	newUser.Password = string(hashPassword)
	users = append(users, newUser)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "User registered successfully",
	})
}
func LoginUser(ctx *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var userOk bool
	if _, ok := userSecrets[loginRequest.Username]; !ok {
		userOk = false
	} else {
		for _, u := range users {
			if u.Username == loginRequest.Username {
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginRequest.Password)); err == nil {
					userOk = true
					break
				}
			}
		}
	}
	if userOk {
		Token, err := GenerateJWT(loginRequest.Username, false)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
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
