package watermill_midlewares

import (
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

func CustomCorrelationIdMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")
		if correlationID == "" {
			correlationID = shortuuid.New()
			msg.Metadata.Set("correlation_id", correlationID)
		}
		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		msg.SetContext(ctx)
		return next(msg)
	}
}
