package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/customsturctures"
	"github.com/Shaieb524/web-clinic.git/helpers"
	"github.com/Shaieb524/web-clinic.git/responses"
	"go.uber.org/zap"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var customLogger *zap.Logger = helpers.CustomLogger()

func JWTauthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken(c)
	}
}

func verifyToken(c *gin.Context) {
	token, ok := getTokenFromHeader(c)
	if !ok {
		customLogger.Error("Error extracting token from header")
		c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
			Data: map[string]interface{}{"error": "Couldn't extract token from header"}})
		return
	}

	email, err := validateToken(token)

	if err != nil {
		customLogger.Error("Error validating token")
		c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
			Data: map[string]interface{}{"error": "Error validating token"}})
		return
	}

	c.Set("useremail", email)
	c.Writer.Header().Set("Authorization", "Bearer "+token)
	c.Next()
}

func getTokenFromHeader(c *gin.Context) (string, bool) {
	authValue := c.GetHeader("Authorization")

	arr := strings.Split(authValue, " ")
	if len(arr) != 2 {
		customLogger.Error("Error inside header array")
		return "", false
	}

	authType := strings.Trim(arr[0], "\n\r\t")
	if strings.ToLower(authType) != strings.ToLower("Bearer") {
		customLogger.Error("Fisrt element is not Bearer")
		return "", false
	}
	token := strings.Trim(arr[1], "\n\t\r")

	return token, true
}

func validateToken(tokenString string) (string, error) {

	var claims customsturctures.AuthClaimers
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			customLogger.Error("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.EnvTokenSecretKey()), nil
	})
	if err != nil {
		customLogger.Error("Problem while parsing token")
		return "__", err
	}
	if !token.Valid {
		customLogger.Error("Token is not valid")
		return "__", errors.New("invalid token")
	}
	email := claims.Email
	return email, nil
}
