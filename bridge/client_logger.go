package bridge

import (
	"github.com/whatsapp-signal-bridge/logger"
)

type clientLogger struct {
	next   Client
	logger logger.Logger
}

func NewClientLogger(client Client) Client {
	return &clientLogger{client, logger.NewLogger("client")}
}

func (l *clientLogger) Publish(messageQueueId string, msg Message) {
	l.logger.Log("publish", msg, "to", messageQueueId)
	l.next.Publish(messageQueueId, msg)
}

func (l *clientLogger) Subscribe(messageQueueId string, callback func(msg Message)) {
	l.logger.Log("new subscription to", messageQueueId)
	l.next.Subscribe(messageQueueId, func(msg Message) {
		l.logger.Log("received message on", messageQueueId, ":", msg)
		callback(msg)
	})
}

func (l *clientLogger) Send(msg Message) {
	l.logger.Log(msg)
	l.next.Send(msg)
}
