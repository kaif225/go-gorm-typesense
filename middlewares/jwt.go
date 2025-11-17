package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"psql-typesense/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() gin.HandlerFunc {

	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenCookie, err := c.Request.Cookie("Bearer")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Authorization token required",
				})
				return
			}
			tokenString = tokenCookie.Value
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		claims := &utils.SignedDetail{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		c.Set("role", claims.Role)
		//c.Set("userID", claims.)

		c.Next()
	}
}
