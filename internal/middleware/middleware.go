package middleware

import (
	"go-api/internal/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(repo repo.UserRepoInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "login before requesting this ressource",
			})
			return
		}
		if !repo.CheckSession(session) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "session expired or invalid",
			})
			return
		}
		ctx.Next()
	}
}
