package handler

import (
	"fmt"
	"time"
	"timeline/config"
	"timeline/internal/handler/middleware"
	"timeline/internal/service"

	"github.com/gin-gonic/gin"
	// "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
	// _ "timeline/docs"
)

type Handler struct {
	s   *service.Service
	cfg *config.Config
}

func NewHandler(service *service.Service, config *config.Config) *Handler {
	return &Handler{
		s:   service,
		cfg: config,
	}
}

func (h *Handler) InitRoute() *gin.Engine {
	router := gin.Default()

	router.Use(gin.LoggerWithFormatter(
		func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}))

	router.Use(
		middleware.Cors(h.cfg),
	)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		h.initUserRoutes(api)
		h.initCompetitionRoutes(api)
		h.initClubRoutes(api)
		h.initSeasonRoutes(api)
		h.initMatchRoutes(api)
	}
	return router
}
