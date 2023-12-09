package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	database "github.com/jyotirmoydotdev/openfy/db"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(ctx *gin.Context) {
	// Structure to hold incoming JSON data
	var newAdmin struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Password  string `json:"password"`
	}
	// Structure to hold admin data for database storage
	var newAdminDatabase database.Admin
	// Bind incoming JSON data to the newAdmin structure
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
	// Convert the email and username to lowercase for case-insensitivity
	newAdmin.Email = strings.ToLower(newAdmin.Email)
	newAdmin.Username = strings.ToLower(newAdmin.Username)
	// Check the length of username
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
	// Check if username is valid (contains only letters and numbers)
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
	// Check if the username or email already exists
	for _, a := range database.Admins {
		if a.Username == newAdmin.Username {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":       "INVALID_INPUT",
				"message":     "Invalid input data",
				"success":     false,
				"field":       "username",
				"description": "username not available",
			})
			return
		}
		if a.Email == newAdmin.Email {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":       "INVALID_INPUT",
				"message":     "Invalid input data",
				"success":     false,
				"field":       "email",
				"description": "email already exists",
			})
			return
		}
	}
	// Combine First and Last name for the database
	newAdminDatabase.Name = newAdmin.FirstName + " " + newAdmin.LastName
	// Hash the password
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
	// Replace the password with the hash
	newAdmin.Password = string(hashPassword)
	// Generate a Secret key for the Admin
	secretKey, err := generateRandomKey()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":       "INTERNAL_SERVER_ERROR",
			"message":     "internal server error",
			"success":     false,
			"field":       "secrerKey",
			"description": "An internal server error occurred during generating secreat key",
		})
		return
	}
	// Keep the secretKey in the adminSecret map
	database.AdminSecrets[newAdmin.Username] = secretKey
	// Generate the Admin ID
	newAdminDatabase.ID = generateAdminID()
	// Set AccountOwner flag based on the number of existing admins
	newAdminDatabase.AccountOwner = len(database.Admins) == 0
	// Copy data from newAdmin to newAdminDatabase
	err = copier.Copy(&newAdminDatabase, &newAdmin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Add the new Admin to the Admins Array
	database.Admins = append(database.Admins, newAdminDatabase)
	// Return success message along with the new admin data
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
	_, adminExists := database.AdminSecrets[loginRequest.Username]
	if !adminExists {
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
	var matchedAdmin database.Admin
	for _, a := range database.Admins {
		if a.Username == loginRequest.Username {
			// Compare the password hash
			if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(loginRequest.Password)); err == nil {
				matchedAdmin = a
				break
			}
		}
	}
	// If a matching admin is found, generate and return a JWT
	if matchedAdmin.ID != "" {
		token, err := GenerateJWT(loginRequest.Username)
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
	database.AdminIDCounter++
	return fmt.Sprintf("A%d", database.AdminIDCounter)
}

func HashAdmin() bool {
	return len(database.Admins) != 0
}
