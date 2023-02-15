package middleware

import (
	"net/http"
	"timeline/config"

	"github.com/gin-gonic/gin"
)

func Cors(cfg *config.Config) gin.HandlerFunc {
    return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", cfg.Client.Url)
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,OPTIONS")
		
		ctx.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method != "OPTIONS" {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusOK)
		}
    }
}
