package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *gin.Context) {
	var newUser struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	// Check if user with the same email already exists
	exists, err := models.UserExistByEmail(dbInstance, newUser.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists, please login",
		})
		return
	}

	// Create a new user hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	userModel := models.NewUserModel(dbInstance)

	newUserDatabase := models.Customer{
		Email:     strings.ToLower(newUser.Email),
		Password:  string(hashPassword),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		DeliveryAddresses: []models.DeliveryAddress{
			{
				Country:   "",
				Address1:  "",
				Apartment: "",
				City:      "",
				Province:  "",
				Zip:       0,
			},
		},
	}
	// Save user to the database
	if err := userModel.Save(&newUserDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	UserSecret, err := generateRandomKey()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while creating secreatKey",
		})
	}
	userID, err := userModel.GetUserID(strings.ToLower(newUser.Email))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	newUserSecret := models.UserSecrets{
		UserID: userID,
		Email:  strings.ToLower(newUser.Email),
		Secret: UserSecret,
	}
	if err := userModel.SaveUserSecret(&newUserSecret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "User registered successfully",
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
	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	loginRequest.Email = strings.ToLower(loginRequest.Email)
	var userOk bool
	if _, err := models.GetUserSecretKeyByEmail(dbInstance, loginRequest.Email); err != nil {
		userOk = false
	} else {
		userHashedPassword, err := models.GetUserHashedPasswordByEmail(dbInstance, loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(loginRequest.Password)); err == nil {
			userOk = true
		}
	}
	if userOk {
		Token, err := GenerateUserJWT(loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server error",
				"message": err.Error(),
			})
			return
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
