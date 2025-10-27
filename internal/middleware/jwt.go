package middleware

import (
	"github.com/Asylann/Auth/lib/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Auth(secret string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")
		if err != nil {
			logger.Errorf("Cant find cookie: %s", err.Error())
			c.IndentedJSON(http.StatusBadRequest, gin.H{"err": "Cant find cookie"})
			c.Abort()
			return
		}
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.Errorf("Token is invalid: %s", err.Error())
			c.IndentedJSON(http.StatusBadRequest, gin.H{"err": "Cant convert jwt token"})
			c.Abort()
			return
		}

		c.Set("User", claims.User)
		c.Set("Role", claims.Role)
		c.Next()
	}
}
