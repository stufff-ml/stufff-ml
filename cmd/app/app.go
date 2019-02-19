package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stufff-ml/stufff-ml/internal/api"
	"github.com/stufff-ml/stufff-ml/internal/callback"
	"github.com/stufff-ml/stufff-ml/internal/jobs"
	"github.com/stufff-ml/stufff-ml/internal/scheduler"
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

	//
	// Public API. Namespace /api/1
	//
	apiNamespace := router.Group(a.APIPrefix)

	// events
	apiNamespace.GET("/events", api.GetEventsEndpoint)
	apiNamespace.POST("/events", api.PostEventsEndpoint)

	// client
	apiNamespace.GET("/client/create", api.ClientCreateEndpoint)

	// model
	apiNamespace.POST("/model/train", api.ModelTrainingEndpoint)
	apiNamespace.GET("/model/predict", api.GetPredictionEndpoint)

	// Admin API. Namespace /_a
	adminNamespace := router.Group(a.AdminAPIPrefix)
	adminNamespace.GET("/init", api.InitEndpoint)

	//
	// internal routes, not part of the official API
	//

	// /_i/1/callback
	callbackNamespace := router.Group(a.CallbackPrefix)
	callbackNamespace.GET("/train", callback.ModelTrainingEndpoint)

	// /_i/1/scheduler
	schedulerNamespace := router.Group(a.SchedulerPrefix)
	schedulerNamespace.GET("/export", scheduler.EventsExportEndpoint)
	schedulerNamespace.GET("/train", scheduler.ModelTrainingEndpoint)

	// /_i/1/jobs
	jobsNamespace := router.Group(a.JobsPrefix)
	jobsNamespace.POST("/export", jobs.EventsExportEndpoint)
	jobsNamespace.POST("/merge", jobs.EventsMergeEndpoint)
	jobsNamespace.POST("/train", jobs.ModelTrainingEndpoint)
	jobsNamespace.POST("/import", jobs.ModelImportEndpoint)

	//
	// default endpoints that are not part of the API namespace
	//
	router.GET("/", api.DefaultEndpoint)
	router.GET("/robots.txt", api.RobotsEndpoint)

	// ready, start taking requests
	http.Handle("/", router)
}
