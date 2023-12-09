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
	database "github.com/jyotirmoydotdev/openfy/db"
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := database.AdminSecrets[username]
	signalToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "Internal Server error", err
	}
	return signalToken, nil
}
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(extractSecretkeyFromToken(token)), nil
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
func extractSecretkeyFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	username, ok := claims["username"].(string)
	if !ok {
		return ""
	}
	return database.AdminSecrets[username]
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
				"error": "Unauthorized ",
			})
			ctx.Abort()
			return
		}
		claims, err := ValidateToken(tokenString)
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
