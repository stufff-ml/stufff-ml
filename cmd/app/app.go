package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stufff-ml/stufff-ml/internal/app"
	"github.com/stufff-ml/stufff-ml/pkg/api"
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
	apiNamespace := router.Group(api.APIBaseURL)
	apiNamespace.GET("/events", app.GetEventsEndpoint)
	apiNamespace.POST("/events", app.PostEventsEndpoint)
	apiNamespace.POST("/predict", app.GetPredictionEndpoint)

	// /_i/1/batch
	batchNamespace := router.Group(api.BatchBaseURL)
	batchNamespace.POST("/predictions", app.PostPredictionsEndpoint)

	// /_i/1/scheduler
	schedulerNamespace := router.Group(api.SchedulerBaseURL)
	schedulerNamespace.GET("/export", app.ScheduleEventsExportEndpoint)

	// /_i/1/jobs
	jobsNamespace := router.Group(api.JobsBaseURL)
	jobsNamespace.POST("/export", app.JobEventsExportEndpoint)
	jobsNamespace.POST("/merge", app.JobEventsMergeEndpoint)

	// namespace /_a
	adminNamespace := router.Group(api.AdminBaseURL)
	adminNamespace.GET("/init", app.InitEndpoint)

	// ready, start taking requests
	http.Handle("/", router)

}
