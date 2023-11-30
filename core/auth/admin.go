package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	database "github.com/jyotirmoydotdev/openfy/Database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(ctx *gin.Context) {
	var newAdmin struct {
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	var newAdminDatabase database.Admin
	if err := ctx.ShouldBindJSON(&newAdmin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	for _, a := range database.Admins {
		if a.Username == newAdmin.Username {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "username not available",
			})
			return
		}
		if a.Email == newAdmin.Email {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "email already exist",
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
	database.AdminSecrets[newAdmin.Username] = secretKey
	newAdmin.Password = string(hashPassword)
	newAdminDatabase.ID = generateAdminID()
	err = copier.Copy(&newAdminDatabase, &newAdmin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	database.Admins = append(database.Admins, newAdminDatabase)
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
	if _, ok := database.AdminSecrets[loginRequest.Username]; !ok {
		adminOk = false
	} else {
		for _, a := range database.Admins {
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
	database.AdminIDCounter++
	return fmt.Sprintf("A%d", database.AdminIDCounter)
}

func HashAdmin() bool {
	return len(database.Admins) != 0
}
