package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	signalToken, err := token.SignedString([]byte(userSecrets[username]))
	if err != nil {
		return "Internal Server error", err
	}
	return signalToken, nil
}
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
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
	return claims, nil
}
func extractUsernameFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	username, ok := claims["sub"].(string)
	if !ok {
		return ""
	}
	return username
}
func extractSecretkeyFromToken(token *jwt.Token) string {
	username := extractUsernameFromToken(token)
	secretkey, ok := userSecrets[username]
	if !ok {
		return ""
	}
	return secretkey
}
