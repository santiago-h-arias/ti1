package middlewares

import (
	"log"
	"net/http"

	services "tinc1/Services"

	"github.com/gin-gonic/gin"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authorizationHeader := c.GetHeader("Authorization")
		tokenString := authorizationHeader[len(BEARER_SCHEMA):]

		token, err := services.NewJWTService().ValidateToken(tokenString)

		if !token.Valid {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
