package cron

import "fmt"

type JobEvent struct {
	ID JobID
}

type JobAddEvent struct {
	JobEvent
}

func (j *JobAddEvent) String() string {
	return fmt.Sprintf("job %v add", j.ID)
}

type JobRemoveEvent struct {
	JobEvent
}

func (j *JobRemoveEvent) String() string {
	return fmt.Sprintf("job %v remove", j.ID)
}

type JobStartEvent struct {
	JobEvent
}

func (j *JobStartEvent) String() string {
	return fmt.Sprintf("job %v start", j.ID)
}

