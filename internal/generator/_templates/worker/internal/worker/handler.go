package worker

import (
	"context"
	"log/slog"

	async "___LIB_PATH___/lib-async"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, msg async.Envelope[Job]) error {
	slog.InfoContext(
		ctx,
		"received worker job",
		"message_uuid", msg.UUID,
		"job_id", msg.Payload.ID,
		"job_type", msg.Payload.Type,
		"payload", string(msg.Payload.Payload),
	)
	return nil
}
