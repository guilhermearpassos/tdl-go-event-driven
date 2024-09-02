package watermill_midlewares

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

func CustomLogginMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		msgs, err := next(msg)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"message_uuid":   msg.UUID,
					"correlation_id": msg.Metadata.Get("correlation_id"),
					"error":          err,
				},
			).Error("Message handling error")
		}
		logrus.WithFields(
			logrus.Fields{
				"message_uuid":   msg.UUID,
				"correlation_id": msg.Metadata.Get("correlation_id"),
			},
		).Info("Handling a message")
		return msgs, err
	}
}
