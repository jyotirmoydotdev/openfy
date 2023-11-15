package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var admins []Admin
var adminSecrets = make(map[string]string)
var adminIDCounter int

func RegisterAdmin(ctx *gin.Context) {
	var newAdmin Admin
	if err := ctx.ShouldBindJSON(&newAdmin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	for _, a := range admins {
		if a.Username == newAdmin.Username {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Admin name not available",
			})
			return
		}
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), bcrypt.DefaultCost)
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
	adminSecrets[newAdmin.Username] = secretKey
	newAdmin.Password = string(hashPassword)
	newAdmin.ID = generateAdminID()
	admins = append(admins, newAdmin)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "User registered successfully",
	})
}

func LoginAdmin(ctx *gin.Context) {
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
	var adminOk bool
	if _, ok := adminSecrets[loginRequest.Username]; !ok {
		adminOk = false
	} else {
		for _, a := range admins {
			if a.Username == loginRequest.Username {
				if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(loginRequest.Password)); err == nil {
					adminOk = true
					break
				}
			}
		}
	}
	if adminOk {
		Token, err := GenerateJWT(loginRequest.Username, true)
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
			"error": "Not a valid Admin",
		})
		return
	}
}

func generateAdminID() string {
	adminIDCounter++
	return fmt.Sprintf("A%d", adminIDCounter)
}

func HashAdmin() bool {
	return len(admins) != 0
}
