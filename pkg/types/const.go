package types

const (
	// Version is the human readable version string of this build
	Version string = "1.0"

	// APIBaseURL is the prefix for all API calls
	APIBaseURL string = "/api/1"
	// BatchBaseURL is the prefix for all batch import/export endpoints
	BatchBaseURL string = "/_i/1/batch"
	// SchedulerBaseURL is the prefix for all scheduller/cron tasks
	SchedulerBaseURL string = "/_i/1/scheduler"
	// JobsBaseURL is the prefix for all scheduled jobs
	JobsBaseURL string = "/_i/1/jobs"
	// AdminBaseURL is the prefix for all admin endpoints
	AdminBaseURL string = "/_a"

	// DefaultExportBucket is the bucket used to export data to
	DefaultExportBucket string = "exports.stufff.review"
	// DefaultModelsBucket is the bucket for model data
	DefaultModelsBucket string = "models.stufff.review"
)
