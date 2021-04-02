package bridge

const (
	SIGNAL_QUEUE   = "signal"
	WHATSAPP_QUEUE = "whatsapp"
)

type Queue interface {
	Publish(messageQueueId string, msg Message)
	Subscribe(callback func(messageQueueId string, msg Message))
}

type queue struct {
	subscriptions []func(messageQueueId string, msg Message)
}

func NewQueue() Queue {
	return NewQueueLogger(&queue{})
}

func (q *queue) Publish(messageQueueId string, msg Message) {
	for _, subscription := range q.subscriptions {
		subscription(messageQueueId, msg)
	}
}

func (q *queue) Subscribe(callback func(messageQueueId string, msg Message)) {
	q.subscriptions = append(q.subscriptions, callback)
}
