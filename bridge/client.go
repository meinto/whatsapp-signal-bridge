package bridge

import "log"

type Client interface {
	Publish(messageQueueId string, msg Message)
	Subscribe(messageQueueId string, callback func(msg Message))
	Send(msg Message)
}

type client struct {
	queue Queue
}

func NewClient(queue Queue) Client {
	return &client{queue}
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

func (c *client) Send(_ Message) {
	log.Println("please implment send method")
}
