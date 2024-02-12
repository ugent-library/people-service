package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/ugent-library/people-service/indexes"
	"github.com/ugent-library/people-service/repositories"
)

type ReindexPeopleArgs struct{}

func (ReindexPeopleArgs) Kind() string { return "reindexPeople" }

func (ReindexPeopleArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{
			ByQueue: true,
		},
	}
}

type ReindexPeopleWorker struct {
	river.WorkerDefaults[ReindexPeopleArgs]
	repo  *repositories.Repo
	index *indexes.Index
}

func NewReindexPeopleWorker(repo *repositories.Repo, index *indexes.Index) *ReindexPeopleWorker {
	return &ReindexPeopleWorker{repo: repo, index: index}
}

func (w *ReindexPeopleWorker) Work(ctx context.Context, job *river.Job[ReindexPeopleArgs]) error {
	return w.index.ReindexPeople(ctx, w.repo.EachPerson)
}
