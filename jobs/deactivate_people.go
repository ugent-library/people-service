package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/ugent-library/people-service/repositories"
)

type DeactivatePeopleArgs struct{}

func (DeactivatePeopleArgs) Kind() string { return "deactivatePeople" }

func (DeactivatePeopleArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByQueue: true,
		},
	}
}

type DeactivatePeopleWorker struct {
	river.WorkerDefaults[DeactivatePeopleArgs]
	repo *repositories.Repo
}

func NewDeactivatePeopleWorker(repo *repositories.Repo) *DeactivatePeopleWorker {
	return &DeactivatePeopleWorker{repo: repo}
}

func (w *DeactivatePeopleWorker) Work(ctx context.Context, job *river.Job[DeactivatePeopleArgs]) error {
	return w.repo.DeactivatePeople(ctx)
}
