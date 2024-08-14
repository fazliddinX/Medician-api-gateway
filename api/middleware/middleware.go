package middleware

import (
	"github.com/casbin/casbin/v2"

	"api-gateway/pkg/token"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, err := token.ExtractClaims(stringToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func Authorization(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		userClaims, ok := claims.(*token.Claims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			ctx.Abort()
		}

		fmt.Println(userClaims.Role, ctx.FullPath(), ctx.Request.Method)
		ok, err := enforcer.Enforce(userClaims.Role, ctx.FullPath(), ctx.Request.Method)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal server error",
				"Err":   err.Error(),
			})
			ctx.Abort()
			return
		}

		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"Error": "Forbidden",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
