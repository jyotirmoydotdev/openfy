package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jyotirmoydotdev/openfy/database"
	"github.com/jyotirmoydotdev/openfy/database/models"
	"gorm.io/gorm"
)

func GenerateJWT(db *gorm.DB, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 336).Unix(),
	})
	secretKey, err := models.GetSecretKeyByUsername(db, username)
	if err != nil {
		if err != nil {
			return "Internal Server error", err
		}
	}
	signalToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "Internal Server error", err
	}
	return signalToken, nil
}
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	dbInstance, err := database.GetCustomerDB()
	if err != nil {
		return nil, err
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(extractSecretkeyFromToken(dbInstance, token)), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	expirationTime, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("something went wrong while extracting 'expirationTime'")
	}
	expiration := time.Unix(int64(expirationTime), 0)
	if !(time.Now().Before(expiration)) {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}
func CheckExpiration(token *jwt.Token) bool {
	expirationTime, ok := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return false
	}
	expiration := time.Unix(int64(expirationTime), 0)
	return time.Now().Before(expiration)
}
func extractSecretkeyFromToken(dbInstance *gorm.DB, token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	username, ok := claims["username"].(string)
	if !ok {
		return ""
	}
	secretKey, err := models.GetSecretKeyByUsername(dbInstance, username)
	if err != nil {
		return ""
	}
	return secretKey
}
func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "Something went wrong", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
func AuthenticateMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized ",
				"reason": "token not found",
			})
			ctx.Abort()
			return
		}
		claims, err := ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized",
				"reason": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
