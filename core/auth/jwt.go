package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	database "github.com/jyotirmoydotdev/openfy/Database"
)

func GenerateJWT(username string, email string, isAdmin bool) (string, error) {
	var role string
	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"email":    email,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	var secretKey string
	if isAdmin {
		secretKey = database.AdminSecrets[username]
	} else {
		secretKey = database.UserSecrets[email]
	}
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
	checkExpiration := CheckExpiration(token)
	if !checkExpiration {
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
	username := extractUsernameFromToken(token)
	Email := extractUsernameFromToken(token)
	var secretkey string
	var isOk bool
	if isAdmin, err := extractIsAdminFromToken(token); err == nil && isAdmin {
		secretkey, isOk = database.AdminSecrets[username]
	} else {
		secretkey, isOk = database.UserSecrets[email]
	}
	if !isOk {
		return ""
	}
	return secretkey
}
func extractUsernameFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	username, ok := claims["username"].(string)
	if !ok {
		return ""
	}
	return username
}
func extractIsAdminFromToken(token *jwt.Token) (bool, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return false, fmt.Errorf("invalid or missing isAdmin field in token claims")
	}
	return role == "admin", nil
}
func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "Something went wrong", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
