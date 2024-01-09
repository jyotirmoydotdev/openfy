package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterCustomer(ctx *gin.Context) {
	var newCustomer struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&newCustomer); err != nil {
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

	// Check if customer with the same email already exists
	exists, err := models.CustomerExistByEmail(dbInstance, newCustomer.Email)
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

	// Create a new customer hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newCustomer.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	customerModel := models.NewCustomerModel(dbInstance)

	newCustomerDatabase := models.Customer{
		Email:     strings.ToLower(newCustomer.Email),
		Password:  string(hashPassword),
		FirstName: newCustomer.FirstName,
		LastName:  newCustomer.LastName,
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
	// Save customer to the database
	if err := customerModel.Save(&newCustomerDatabase); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	CustomerSecret, err := generateRandomKey()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while creating secreatKey",
		})
	}
	customerID, err := customerModel.GetCustomerID(strings.ToLower(newCustomer.Email))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	newCustomerSecret := models.CustomerSecrets{
		CustomerID: customerID,
		Email:      strings.ToLower(newCustomer.Email),
		Secret:     CustomerSecret,
	}
	if err := customerModel.SaveCustomerSecret(&newCustomerSecret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Customer registered successfully",
	})
}

func LoginCustomer(ctx *gin.Context) {
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
	var customerOk bool
	if _, err := models.GetCustomerSecretKeyByEmail(dbInstance, loginRequest.Email); err != nil {
		customerOk = false
	} else {
		customerHashedPassword, err := models.GetCustomerHashedPasswordByEmail(dbInstance, loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(customerHashedPassword), []byte(loginRequest.Password)); err == nil {
			customerOk = true
		}
	}
	if customerOk {
		Token, err := GenerateCustomerJWT(loginRequest.Email)
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
			"error": "Not a valid customer",
		})
		return
	}
}

func UpdateCustomer(ctx *gin.Context) {
}
