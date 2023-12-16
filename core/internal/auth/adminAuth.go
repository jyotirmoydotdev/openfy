package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
	"golang.org/x/crypto/bcrypt"
)

func SignupAdmin(ctx *gin.Context) {
	var newAdmin struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Password  string `json:"password"`
	}

	var newAdminDatabase models.Admin

	if err := ctx.ShouldBindJSON(&newAdmin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_JSON",
			"message":     "Invalid JSON request",
			"success":     false,
			"field":       "",
			"description": "The request body must be a valid JSON object",
		})
		return
	}

	newAdmin.Email = strings.ToLower(newAdmin.Email)
	newAdmin.Username = strings.ToLower(newAdmin.Username)

	if len(newAdmin.Username) < 4 || len(newAdmin.Username) > 16 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "username",
			"description": "username must be minimum of 3 and maximun of 16 character",
		})
		return
	}

	for _, c := range newAdmin.Username {
		if (97 <= c && c <= 122) || (48 <= c && c <= 57) {
			continue
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":       "INVALID_INPUT",
				"message":     "Invalid input data",
				"success":     false,
				"field":       "username",
				"description": "username can only contain letters and numbers",
			})
			return
		}
	}

	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// Check if a username exist
	if usernameExist, err := models.AdminExistByUsername(dbInstance, newAdmin.Username); err != nil && !usernameExist {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "username",
			"description": "username not available",
		})
		return
	}

	// Check if a email exist
	if emailExist, err := models.AdminExistByEmail(dbInstance, newAdmin.Email); err != nil && !emailExist {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "email",
			"description": "email already exists",
		})
		return
	}

	newAdminDatabase.Name = newAdmin.FirstName + " " + newAdmin.LastName

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "INTERNAL_SERVER_ERROR",
			"message":     "internal server error",
			"success":     false,
			"field":       "password",
			"description": "An internal server error occurred during password hashing",
		})
		return
	}

	newAdmin.Password = string(hashPassword)

	secretKey, err := generateRandomKey()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "INTERNAL_SERVER_ERROR",
			"message":     "internal server error",
			"success":     false,
			"field":       "secrertKey",
			"description": "An internal server error occurred during generating secreat key",
		})
		return
	}
	newAdminSecret := models.AdminSecrets{
		Username: newAdmin.Username,
		Secret:   secretKey,
	}
	adminModel := models.NewAdminModel(dbInstance)

	newAdminDatabase.ID = generateAdminID()

	// Check if the admin table id empty
	AccountOwner, err := models.CheckAdminTableIsEmpty(dbInstance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	newAdminDatabase.AccountOwner = AccountOwner

	err = copier.Copy(&newAdminDatabase, &newAdmin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := adminModel.Save(&newAdminDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if err := adminModel.SaveAdminSecret(&newAdminSecret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":       "",
		"message":     "",
		"success":     true,
		"field":       "",
		"description": "Admin Registered Successfully",
	})
}

func LoginAdmin(ctx *gin.Context) {
	// Structure to hold incoming JSON data
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_JSON",
			"message":     "Invalid JSON request",
			"success":     false,
			"field":       "",
			"description": "The request body must be a valid JSON object",
		})
		return
	}
	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	if check, err := models.AdminExistByUsername(dbInstance, loginRequest.Username); err != nil || check {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_ADMIN",
			"message":     "Invalid Admin",
			"success":     false,
			"field":       "",
			"description": "The admin doesn't exit",
		})
		return
	}
	// Find the admin by username
	adminHashedPassword, err := models.GetAdminHashedPasswordByUsername(dbInstance, loginRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	// Compare the password hash
	// If a match admin is found, generate and return a JWT
	if err := bcrypt.CompareHashAndPassword([]byte(adminHashedPassword), []byte(loginRequest.Password)); err == nil {
		token, err := GenerateJWT(dbInstance, loginRequest.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":       "INTERNAL_SERVER_ERROR",
				"message":     "Internal Server error",
				"success":     false,
				"field":       "JWT Token",
				"description": "An internal server error occurred during generating token",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"error":       "",
			"message":     "",
			"success":     true,
			"field":       "JWT Token",
			"description": "Token generated successfully",
			"token":       token,
		})
		return
	}
	// No matching admin found
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error":       "INVALID_ADMIN",
		"message":     "invalid admin",
		"success":     false,
		"field":       "",
		"description": "No admin found",
	})
}

func generateAdminID() string {
	models.AdminIDCounter++
	return fmt.Sprintf("A%d", models.AdminIDCounter)
}

func HashAdmin() (bool, error) {
	dbInstance, err := db.GetDB()
	if err != nil {
		return false, err
	}
	isEmpty, err := models.CheckAdminTableIsEmpty(dbInstance)
	if err != nil {
		return false, err
	}
	return isEmpty, nil
}
