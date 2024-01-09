package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jyotirmoydotdev/openfy/db"
	"github.com/jyotirmoydotdev/openfy/db/models"
	"golang.org/x/crypto/bcrypt"
)

func GenerateCustomerJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 336).Unix(), // Token Valid for 14 Days
	})
	dbInstance, err := db.GetDB()
	if err != nil {
		return "", err
	}
	secretKey, err := models.GetCustomerSecretKeyByEmail(dbInstance, email)
	if err != nil {
		return "", err
	}
	signToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	parts := strings.Split(signToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT token format")
	}
	payload := parts[1]
	if len(payload) > 72 {
		payload = payload[len(payload)-72:]
	}
	hashToken, err := bcrypt.GenerateFromPassword([]byte(payload), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}

	customerToken := models.CustomerToken{
		Email:             email,
		Token:             string(hashToken),
		LastUsed:          time.Now(),
		TokenExpiry:       time.Now().Add(time.Hour * 24),
		IsActive:          true,
		IPAddresses:       "",
		CustomerAgent:     "",
		DeviceInformation: "",
		RevocationReason:  "",
	}
	if exist, err := models.CheckEmailExist(dbInstance, email); err != nil {
		return "", err
	} else if !exist {
		if err := models.SaveToken(dbInstance, &customerToken); err != nil {
			return "", err
		}
	} else {
		if err := models.UpdateToken(dbInstance, &customerToken); err != nil {
			return "", err
		}
	}
	return signToken, nil
}
func ValidateCustomerToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(extractCustomerSecretKeyFromToken(t)), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("something went wrong at 'claims'")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("something went wrong while extracting `email`")
	}
	dbInstance, err := db.GetDB()
	if err != nil {
		return nil, err
	}
	// Check if the token exist in the database
	databaseToken, err := models.GetTokenByEmail(dbInstance, email)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(tokenString, ".")
	payload := parts[1]
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT token format")
	}
	if len(payload) > 72 {
		payload = payload[len(payload)-72:]
	}
	err = bcrypt.CompareHashAndPassword([]byte(databaseToken), []byte(payload))
	hashTokenCheck := err == nil

	if !hashTokenCheck {
		return nil, fmt.Errorf("not a valid token, not exist in the database")
	}
	expirationTime, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("something went wrong while extracting `expirationTime`")
	}
	expiration := time.Unix(int64(expirationTime), 0)
	if !(time.Now().Before(expiration)) {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}
func extractCustomerSecretKeyFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	email, ok := claims["email"].(string)
	if !ok {
		return ""
	}
	dbInstance, err := db.GetDB()
	if err != nil {
		return ""
	}
	secretKey, err := models.GetCustomerSecretKeyByEmail(dbInstance, email)
	if err != nil {
		return ""
	}
	return secretKey
}
func AuthenticateCustomerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized ",
			})
			ctx.Abort()
			return
		}
		claims, err := ValidateCustomerToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
