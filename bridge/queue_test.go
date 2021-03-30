package bridge

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue()

	TEST_QUEUE := "TEST_QUEUE"
	testMessage := PlainMessage().
		SetChatID("mock")
	subscriptionsCalled := []bool{false, false, false}

	testSubscription := func(index int) func(messageQueueId string, message Message) {
		return func(messageQueueId string, message Message) {
			subscriptionsCalled[index] = true
			if messageQueueId != TEST_QUEUE {
				t.Error("messageQueueId not equal")
			}
			if message.ChatID() != "mock" {
				t.Errorf("message not equal")
			}
		}
	}

	for index := range subscriptionsCalled {
		queue.Subscribe(testSubscription(index))
	}
	queue.Publish(TEST_QUEUE, testMessage)

	for _, subscriptionCall := range subscriptionsCalled {
		if subscriptionCall == false {
			t.Error("subscription not called")
		}
	}
}
