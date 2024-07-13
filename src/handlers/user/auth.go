package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c * gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 401,
				"message": "user unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 401,
				"message": "Invalid Token",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Subject)
		c.Next()
	}
}

func AuthMiddlewareAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 401,
				"message":    "User unauthorized",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code": 401,
				"message":    "Invalid token",
			})
			c.Abort()
			return
		}

		if claims.Subject != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error_code": 403,
				"message":    "Access forbidden: only admin can access",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Subject)
		c.Next()
	}
}
