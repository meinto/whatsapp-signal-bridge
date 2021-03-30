package bridge

import "testing"

func TestClient_Publish(t *testing.T) {
	q := NewMockQueue()
	q.Reset()
	client := NewClient(q)

	MOCK_QUEUE_ID := "MOCK_QUEUE"
	mockMessage := PlainMessage().SetChatID("mock")

	client.Publish(MOCK_QUEUE_ID, mockMessage)
	if q.(*mockQueue).PublishCalls[0].MessageQueueId != MOCK_QUEUE_ID {
		t.Error("message queue id not equal")
	}
	if q.(*mockQueue).PublishCalls[0].Msg.ChatID() != "mock" {
		t.Error("message queue id not equal")
	}
}

func TestClient_Subscribe(t *testing.T) {
	q := NewMockQueue()
	q.Reset()
	client := NewClient(q)

	subscriptionCalled := false
	MOCK_QUEUE_ID := "MOCK_QUEUE"
	mockMessage := PlainMessage().SetChatID("mock")

	client.Subscribe(MOCK_QUEUE_ID, func(msg Message) {
		subscriptionCalled = true
		if msg.ChatID() != mockMessage.ChatID() {
			t.Fail()
		}
	})

	q.(*mockQueue).SubscriptionCalls[0](MOCK_QUEUE_ID, mockMessage)
	if subscriptionCalled == false {
		t.Error("Subscription was not called")
	}

	subscriptionCalled = false
	q.(*mockQueue).SubscriptionCalls[0]("TEST", mockMessage)
	if subscriptionCalled != false {
		t.Error("Subscription should not have beend called")
	}
}
