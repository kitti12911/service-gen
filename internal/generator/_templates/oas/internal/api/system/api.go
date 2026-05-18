package system

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	humautil "___LIB_PATH___/lib-util/v3/huma"

	"___MODULE___/internal/api"
)

func Register(h huma.API, deps api.Deps) {
	huma.Get(h, "/health", func(_ context.Context, _ *struct{}) (*HealthOutput, error) {
		out := &HealthOutput{}
		out.Body.Status = "ok"
		out.Body.Service = deps.ServiceName
		return out, nil
	}, humautil.WithTag(api.TagSystem))
}
