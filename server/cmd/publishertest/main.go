package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/logger"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging"
	"github.com/mrbananaaa/minisocial/internal/platform/messaging/jetstream"
	postDomain "github.com/mrbananaaa/minisocial/internal/post/domain"
	userDomain "github.com/mrbananaaa/minisocial/internal/user/domain"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx := context.Background()

	log := logger.New(&logger.Config{
		Level: "debug",
	})

	broker := jetstream.New(nats.DefaultURL)
	if err := broker.Connect(ctx); err != nil {
		log.Error("failed to connect to JetStream",
			"err", err.Error(),
		)
		return
	}
	defer broker.Close(ctx)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup

	for {
		select {
		case <-ticker.C:
			wg.Go(func() {
				if err := createUser(ctx, broker); err != nil {
					log.Error("failed to create user event",
						"err", err.Error(),
					)
					return
				}
			})

			wg.Go(func() {
				if err := createPost(ctx, broker); err != nil {
					log.Error("failed to create post event",
						"err", err.Error(),
					)
					return
				}
			})

			wg.Wait()
			log.Info("event published!")

		case <-ctx.Done():
			return
		}
	}
}

func createUser(ctx context.Context, broker messaging.Broker) error {
	postCreated := userDomain.UserCreated{
		UserID: uuid.New(),
		Email:  gofakeit.Email(),
	}

	p, err := json.Marshal(postCreated)
	if err != nil {
		return fmt.Errorf("failed to marshal")
	}

	if err := broker.Publish(ctx, "user.created", p); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

func createPost(ctx context.Context, broker messaging.Broker) error {
	postCreated := postDomain.PostCreated{
		PostID:   uuid.New(),
		AuthorID: uuid.New(),
		Title:    "this is published title from event loop",
	}

	p, err := json.Marshal(postCreated)
	if err != nil {
		return fmt.Errorf("failed to marshal")
	}

	if err := broker.Publish(ctx, "post.created", p); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}
