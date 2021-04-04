package bridge

import (
	"github.com/whatsapp-signal-bridge/logger"
)

type clientLogger struct {
	next   Client
	logger logger.Logger
}

func NewClientLogger(client Client) Client {
	return &clientLogger{client, logger.NewLogger("client", logger.LOG_LEVEL_DEBUG)}
}

func (l *clientLogger) Publish(messageQueueId string, msg Message) {
	l.logger.LogDebug("publish", msg, "to", messageQueueId)
	l.next.Publish(messageQueueId, msg)
}

func (l *clientLogger) Subscribe(messageQueueId string, callback func(msg Message)) {
	l.logger.LogDebug("new subscription to", messageQueueId)
	l.next.Subscribe(messageQueueId, func(msg Message) {
		l.logger.LogDebug("received message on", messageQueueId, ":", msg)
		callback(msg)
	})
}

func (l *clientLogger) Send(msg Message) {
	l.logger.LogDebug(msg)
	l.next.Send(msg)
}
