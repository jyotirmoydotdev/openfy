package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	database "github.com/jyotirmoydotdev/openfy/db/repositories"
)

func GenerateUserJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := database.UserSecrets[email]
	signToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "Internal server error", err
	}
	return signToken, nil
}
func ValidateUserToken(tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(extractUserSecretKeyFromToken(t)), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("something went wrong at 'claims'")
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
func extractUserSecretKeyFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	email, ok := claims["email"].(string)
	if !ok {
		return ""
	}
	return database.UserSecrets[email]
}
func AuthenticateUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized ",
			})
			ctx.Abort()
			return
		}
		claims, err := ValidateUserToken(tokenString)
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
