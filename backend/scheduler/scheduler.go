package scheduler

import (
	"context"
	"log"

	"github.com/koo-arch/adjusta-backend/cache"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	client *ent.Client
	cron  *cron.Cron
	keyManager *auth.KeyManager
}

func NewScheduler(client *ent.Client, cache *cache.Cache) *Scheduler {
	return &Scheduler{
		client: client,
		cron: cron.New(),
		keyManager: auth.NewKeyManager(client, cache),
	}
}

func (s *Scheduler) SetupJobs(ctx context.Context) {
	_, err := s.cron.AddFunc("@daily", func() {
		// キーの生成
		keyManager := s.keyManager
		err := keyManager.GenerateJWTKey(ctx, "access")
		if err != nil {
			log.Printf("failed to generate JWT key: %v", err)
		}
		err = keyManager.GenerateJWTKey(ctx, "refresh")
		if err != nil {
			log.Printf("failed to generate JWT key: %v", err)
		}

		// キーの削除
		err = keyManager.DeleteJWTKeys(ctx, "access")
		if err != nil {
			log.Printf("failed to delete access JWT key: %v", err)
		}
		err = keyManager.DeleteJWTKeys(ctx, "refresh")
		if err != nil {
			log.Printf("failed to delete refresh JWT key: %v", err)
		}

	})
	if err != nil {
		log.Printf("failed to add cron job: %v", err)
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("scheduler stopped")
}