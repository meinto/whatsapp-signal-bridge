package bridge

import "log"

type Client interface {
	Publish(messageQueueId string, msg Message)
	Subscribe(messageQueueId string, callback func(msg Message))
	GetName() string
	Send(msg Message)
}

type client struct {
	queue Queue
	name  string
}

func NewClient(queue Queue, name string) Client {
	return NewClientLogger(&client{queue, name})
}

func (c *client) Publish(messageQueueId string, msg Message) {
	c.queue.Publish(messageQueueId, msg)
}

func (c *client) Subscribe(queueId string, callback func(msg Message)) {
	c.queue.Subscribe(func(messageQueueId string, msg Message) {
		if messageQueueId == queueId {
			callback(msg)
		}
	})
}

func (c *client) GetName() string {
	return c.name
}

func (c *client) Send(_ Message) {
	log.Println("please implment send method")
}
