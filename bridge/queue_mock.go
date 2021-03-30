package bridge

type MockQueue interface {
	Queue
	Reset()
}

type publishCall struct {
	MessageQueueId string
	Msg            Message
}
type subscriptionCall func(messageQueueId string, msg Message)
type mockQueue struct {
	PublishCalls      []publishCall
	SubscriptionCalls []subscriptionCall
}

func NewMockQueue() MockQueue {
	return &mockQueue{}
}

func (q *mockQueue) Reset() {
	q.PublishCalls = []publishCall{}
	q.SubscriptionCalls = []subscriptionCall{}
}
func (q *mockQueue) Publish(messageQueueId string, msg Message) {
	q.PublishCalls = append(q.PublishCalls, publishCall{MessageQueueId: messageQueueId, Msg: msg})
}
func (q *mockQueue) Subscribe(callback func(messageQueueId string, msg Message)) {
	q.SubscriptionCalls = append(q.SubscriptionCalls, callback)
}
