package jobs

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"

	"github.com/majordomusio/commons/pkg/gae/logger"
)

const (
	datastoreJobs string = "JOBS"
)

// Job holds information about a (periodic) job
type Job struct {
	Name    string
	Count   int
	LastRun int64
}

// LastRun retrieves the timestamp a job has last run
func LastRun(ctx context.Context, name string) int64 {
	key := datastore.NewKey(ctx, datastoreJobs, name, 0, nil)
	var job Job

	err := datastore.Get(ctx, key, &job)
	if err != nil {
		return 0
	}
	return job.LastRun
}

// UpdateLastRun updates the job with a timestamp when it has last run
func UpdateLastRun(ctx context.Context, name string, ts int64) error {
	key := datastore.NewKey(ctx, datastoreJobs, name, 0, nil)
	var job Job

	err := datastore.Get(ctx, key, &job)
	if err != nil {
		job = Job{
			name,
			1,
			ts,
		}
	} else {
		job.Count = job.Count + 1
		job.LastRun = ts
	}
	_, err = datastore.Put(ctx, key, &job)
	return err
}

// ScheduleJob is a shorthand to create a background job
func ScheduleJob(ctx context.Context, queue, request string) error {
	t := taskqueue.NewPOSTTask(request, nil)
	_, err := taskqueue.Add(ctx, t, queue)
	if err != nil {
		logger.Error(ctx, "jobs.schedule", err.Error())
	}

	return err
}
