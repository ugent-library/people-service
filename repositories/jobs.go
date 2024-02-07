package repositories

import (
	"context"

	"github.com/riverqueue/river"
)

type deactivatePeopleArgs struct{}

func (deactivatePeopleArgs) Kind() string { return "deactivatePeople" }

type deactivatePeopleWorker struct {
	river.WorkerDefaults[deactivatePeopleArgs]
	repo *Repo
}

func newDeactivatePeopleWorker(repo *Repo) *deactivatePeopleWorker {
	return &deactivatePeopleWorker{repo: repo}
}

func (w *deactivatePeopleWorker) Work(ctx context.Context, job *river.Job[deactivatePeopleArgs]) error {
	return w.repo.DeactivatePeople(ctx)
}
