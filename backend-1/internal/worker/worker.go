package worker

import (
	"context"
	"log"
	"time"

	"matcher/internal/matcher"
	"matcher/internal/repo"
)

type Worker struct {
	Repo     *repo.Repo
	Matcher  *matcher.Service
	Interval time.Duration
}

func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	ids, err := w.Repo.FindRoutesToMatch(ctx, 24*time.Hour)
	if err != nil {
		log.Printf("worker FindRoutesToMatch: %v", err)
		return
	}
	for _, id := range ids {
		if _, err := w.Matcher.MatchRoute(ctx, id); err != nil {
			log.Printf("MatchRoute(%d): %v", id, err)
		}
	}
}
