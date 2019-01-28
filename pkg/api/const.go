package api

const (
	// Version is the human readable version string of this build
	Version string = "1.0"

	// APIBaseURL is the base url of the API
	APIBaseURL string = "https://api.stufff.review"

	// APIPrefix is the prefix for all API calls
	APIPrefix string = "/api/1"
	// AdminAPIPrefix is the prefix for all admin endpoints
	AdminAPIPrefix string = "/_a"

	// BatchPrefix is the prefix for all batch import/export endpoints
	BatchPrefix string = "/_i/1/batch"
	// SchedulerPrefix is the prefix for all scheduller/cron tasks
	SchedulerPrefix string = "/_i/1/scheduler"
	// JobsPrefix is the prefix for all scheduled jobs
	JobsPrefix string = "/_i/1/jobs"
	// CallbackPrefix is the prefix for all callbacks
	CallbackPrefix string = "/_i/1/callback"

	// DefaultExportBucket is the bucket used to export data to
	DefaultExportBucket string = "exports.stufff.review"
	// DefaultModelsBucket is the bucket for model data
	DefaultModelsBucket string = "models.stufff.review"
	// DefaultResourcesBucket is the bucket for binary data etc.
	DefaultResourcesBucket string = "resources.stufff.review"
)
