package database

import (
	"context"

	"___MODULE___/internal/config"

	orm "github.com/kitti12911/lib-orm/v3"
)

func New(ctx context.Context, cfg *config.Config) (*orm.DB, error) {
	return orm.New(
		ctx,
		cfg.Database,
		orm.WithApplicationName(cfg.Service.Name),
		orm.WithTracing(cfg.Tracing.Enabled),
	)
}
