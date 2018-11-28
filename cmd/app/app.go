package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stufff-ml/stufff-ml/internal/app"
	"github.com/stufff-ml/stufff-ml/pkg/types"
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

	// FIXME: DEBUGGING ONLY
	router.GET("/debug", app.DebugEndpoint)
	router.GET("/seed", app.SeedEndpoint)
	router.GET("/migrate", app.MigrateEndpoint)

	// namespace /api/1
	apiNamespace := router.Group(types.APINamespace)
	apiNamespace.GET("/events", app.GetEventEndpoint)
	apiNamespace.POST("/events", app.PostEventEndpoint)

	// ready, start taking requests
	http.Handle("/", router)

}
