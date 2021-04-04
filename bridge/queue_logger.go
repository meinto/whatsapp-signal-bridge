package bridge

import (
	"github.com/whatsapp-signal-bridge/logger"
)

type queueLogger struct {
	next   Queue
	logger logger.Logger
}

func NewQueueLogger(queue Queue) Queue {
	return &queueLogger{queue, logger.NewLogger("queue", logger.LOG_LEVEL_DEBUG)}
}

func (l *queueLogger) Publish(messageQueueId string, msg Message) {
	l.logger.LogDebug("publish", msg, "to", messageQueueId)
	l.next.Publish(messageQueueId, msg)
}

func (l *queueLogger) Subscribe(callback func(messageQueueId string, msg Message)) {
	l.logger.LogDebug("new subscription")
	l.next.Subscribe(func(messageQueueId string, msg Message) {
		l.logger.LogDebug("received message on", messageQueueId, ":", msg)
		callback(messageQueueId, msg)
	})
}
