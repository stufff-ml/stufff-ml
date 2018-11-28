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

	// FIXME: DEBUGGING ONLY
	router.GET("/debug", app.DebugEndpoint)
	router.GET("/seed", app.SeedEndpoint)
	router.GET("/migrate", app.MigrateEndpoint)

	/*
		// API namespace
		apiNS := router.Group(types.APIBaseURL)
		apiNS.GET("/auth", app.AuthEndpoint)
		apiNS.GET("/w", app.GetWorkspacesEndpoint)
		apiNS.GET("/c/:team", app.GetChannelsEndpoint)
		apiNS.GET("/m/:team/:channel", app.GetMessagesEndpoint)
		// health checks
		apiNS.GET("/ready", app.CheckReadyEndpoint)
		apiNS.GET("/alive", app.CheckAliveEndpoint)

		// group of internal endpoints for scheduling
		scheduleNS := router.Group(types.SchedulerBaseURL)
		scheduleNS.GET("/workspace", scheduler.UpdateWorkspaces)
		scheduleNS.GET("/msgs", scheduler.CollectMessages)

		// group of internal endpoints for jobs
		jobsNS := router.Group(types.JobsBaseURL)
		jobsNS.POST("/workspace", jobs.UpdateWorkspaceJob)
		jobsNS.POST("/users", jobs.UpdateUsersJob)
		jobsNS.POST("/channels", jobs.UpdateChannelsJob)
		jobsNS.POST("/msgs", jobs.CollectMessagesJob)
	*/

	// ready, start taking requests
	http.Handle("/", router)

}
