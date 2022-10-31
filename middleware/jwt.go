package middleware

import (
	"api-payment/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			userID := claims["user_id"].(string)
			c.Set("currentUser", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
	}
}