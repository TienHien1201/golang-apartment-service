package xesmigration

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
)

type Migration interface {
	Version() int
	Up(ctx context.Context, client *elasticsearch.Client) error
	Down(ctx context.Context, client *elasticsearch.Client) error
}

type MigrationInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	AppliedAt time.Time `json:"applied_at"`
}
