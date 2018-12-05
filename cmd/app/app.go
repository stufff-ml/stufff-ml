package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stufff-ml/stufff-ml/internal/app"
)

func init() {

	// configure GIN first
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	// a new router
	router := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// default endpoints that are not part of the API namespace
	router.GET("/", app.DefaultEndpoint)
	router.GET("/robots.txt", app.RobotsEndpoint)

	// namespace /api/1
	apiNamespace := router.Group("/api/1")
	apiNamespace.GET("/events", app.GetEventsEndpoint)
	apiNamespace.POST("/events", app.PostEventsEndpoint)
	apiNamespace.POST("/predict", app.GetPredictionEndpoint)

	// internal/integration namespace /_i/1
	internalNamespace := router.Group("/_i/1")
	internalNamespace.POST("/prediction", app.PostPredictionsEndpoint)

	// namespace /_admin
	adminNamespace := router.Group("/_a")
	adminNamespace.GET("/init", app.InitEndpoint)

	// ready, start taking requests
	http.Handle("/", router)

}
