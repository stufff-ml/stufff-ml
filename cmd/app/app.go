package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stufff-ml/stufff-ml/internal/api"
	a "github.com/stufff-ml/stufff-ml/pkg/api"
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
	router.GET("/", api.DefaultEndpoint)
	router.GET("/robots.txt", api.RobotsEndpoint)

	// namespace /api/1
	apiNamespace := router.Group(a.APIBaseURL)
	apiNamespace.GET("/events", api.GetEventsEndpoint)
	apiNamespace.POST("/events", api.PostEventsEndpoint)
	apiNamespace.POST("/predict", api.GetPredictionEndpoint)

	// /_i/1/batch
	batchNamespace := router.Group(a.BatchBaseURL)
	batchNamespace.POST("/predictions", api.PostPredictionsEndpoint)

	// /_i/1/scheduler
	schedulerNamespace := router.Group(a.SchedulerBaseURL)
	schedulerNamespace.GET("/export", api.ScheduleEventsExportEndpoint)

	// /_i/1/jobs
	jobsNamespace := router.Group(a.JobsBaseURL)
	jobsNamespace.POST("/export", api.JobEventsExportEndpoint)
	jobsNamespace.POST("/merge", api.JobEventsMergeEndpoint)

	// namespace /_a
	adminNamespace := router.Group(a.AdminBaseURL)
	adminNamespace.GET("/init", api.InitEndpoint)
	adminNamespace.GET("/client.create", api.ClientCreateEndpoint)

	// ready, start taking requests
	http.Handle("/", router)

}
