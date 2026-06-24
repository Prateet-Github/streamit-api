package middlewares

import (
	"strings"

	"github.com/Prateet-Github/streamit-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "missing token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		claims, err := utils.VerifyToken(
			tokenString,
			secret,
		)

		if err != nil {
			c.JSON(401, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		c.Set(
			"userId",
			claims["userId"],
		)

		c.Next()
	}
}