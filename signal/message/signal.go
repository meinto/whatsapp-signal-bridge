package message

type SignalCLIMessage struct {
	Envelope struct {
		Source         string `json:"source,omitempty"`
		SourceDevice   int    `json:"sourceDevice,omitempty"`
		Timestamp      int    `json:"timestamp,omitempty"`
		ReceiptMessage struct {
			When       int   `json:"when,omitempty"`
			IsDelivery bool  `json:"isDelivery,omitempty"`
			IsRead     bool  `json:"isRead,omitempty"`
			Timestamps []int `json:"timestamps,omitempty"`
		} `json:"receiptMessage,omitempty"`
		DataMessage *struct {
			Timestamp        int    `json:"timestamp,omitempty"`
			Message          string `json:"message,omitempty"`
			ExpiresInSeconds int    `json:"expiresInSeconds,omitempty"`
			ViewOnce         bool   `json:"viewOnce,omitempty"`
			Attachments      []struct {
				ContentType string  `json:"contentType,omitempty"`
				Filename    *string `json:"filename,omitempty"`
				ID          string  `json:"id,omitempty"`
				Size        int     `json:"size,omitempty"`
			} `json:"attachments,omitempty"`
			// Mentions ???
			Quote *struct {
				ID          int    `json:"id,omitempty"`
				Author      string `json:"author,omitempty"`
				Text        string `json:"text,omitempty"`
				Attachments []struct {
					ContentType string  `json:"contentType,omitempty"`
					Filename    *string `json:"filename,omitempty"`
					ID          string  `json:"id,omitempty"`
					Size        int     `json:"size,omitempty"`
				} `json:"attachments,omitempty"`
			} `json:"quote,omitempty"`
		} `json:"dataMessage,omitempty"`
	} `json:"envelope,omitempty"`
}
