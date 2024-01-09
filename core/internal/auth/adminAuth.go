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
	"gorm.io/gorm"
)

type DBInstance struct {
	DB *gorm.DB
}

func SignupStaffMember(ctx *gin.Context) {
	var newStaffMember struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Password  string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&newStaffMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_JSON",
			"message":     "Invalid JSON request",
			"success":     false,
			"field":       "",
			"description": "The request body must be a valid JSON object",
		})
		return
	}

	newStaffMember.Email = strings.ToLower(newStaffMember.Email)
	newStaffMember.Username = strings.ToLower(newStaffMember.Username)

	// Validate the username
	// Check if the username is atleast 4 character or maximum 16 character
	if len(newStaffMember.Username) < 4 || len(newStaffMember.Username) > 16 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "username",
			"description": "username must be minimum of 3 and maximun of 16 character",
		})
		return
	}

	// Check if the username is small character and number
	for _, c := range newStaffMember.Username {
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

	// Get the Database instance to save the data
	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	// Check if a username exist in the database
	if usernameExist, err := models.StaffMemberExistByUsername(dbInstance, newStaffMember.Username); err != nil && !usernameExist {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "username",
			"description": "username not available",
		})
		return
	}

	// Check if a email exist in the database
	if emailExist, err := models.StaffMemberExistByEmail(dbInstance, newStaffMember.Email); err != nil && !emailExist {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_INPUT",
			"message":     "Invalid input data",
			"success":     false,
			"field":       "email",
			"description": "email already exists",
		})
		return
	}

	var newStaffMemberDatabase models.StaffMember

	newStaffMemberDatabase.Name = newStaffMember.FirstName + " " + newStaffMember.LastName

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newStaffMember.Password), bcrypt.DefaultCost)
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

	newStaffMember.Password = string(hashPassword)

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
	staffMemberModel := models.NewStaffMemberModel(dbInstance)

	// Check if the staffMember table id empty
	AccountOwner, err := models.CheckStaffMemberTableIsEmpty(dbInstance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	newStaffMemberDatabase.AccountOwner = AccountOwner

	err = copier.Copy(&newStaffMemberDatabase, &newStaffMember)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := staffMemberModel.Save(&newStaffMemberDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	newStaffMemberId, err := staffMemberModel.GetStaffMemberID(newStaffMember.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	newStaffMemberSecret := models.StaffMemberSecrets{
		StaffMemberID: newStaffMemberId,
		Username:      newStaffMember.Username,
		Secret:        secretKey,
	}

	if err := staffMemberModel.SaveStaffMemberSecret(&newStaffMemberSecret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":     true,
		"description": "StaffMember Registered Successfully",
	})
}

func LoginStaffMember(ctx *gin.Context) {
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
	loginRequest.Username = strings.ToLower(loginRequest.Username)

	dbInstance, err := db.GetDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	if check, err := models.StaffMemberExistByUsername(dbInstance, loginRequest.Username); err != nil || !check {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "INVALID_ADMIN",
			"message":     "Invalid StaffMember",
			"success":     false,
			"field":       "",
			"description": "The staffMember doesn't exit",
		})
		return
	}
	// Find the staffMember by username
	staffMemberHashedPassword, err := models.GetStaffMemberHashedPasswordByUsername(dbInstance, loginRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	// Compare the password hash
	// If a match staffMember is found, generate and return a JWT
	if err := bcrypt.CompareHashAndPassword([]byte(staffMemberHashedPassword), []byte(loginRequest.Password)); err == nil {
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
	// No matching staffMember found
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error":       "INVALID_ADMIN",
		"message":     "invalid staffMember",
		"success":     false,
		"field":       "x",
		"description": "No staffMember found",
	})
}

func HashStaffMember() (bool, error) {
	dbInstance, err := db.GetDB()
	if err != nil {
		return false, err
	}
	isEmpty, err := models.CheckStaffMemberTableIsEmpty(dbInstance)
	if err != nil {
		return false, err
	}
	return isEmpty, nil
}
