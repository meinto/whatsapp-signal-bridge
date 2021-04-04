package whatsapp

import (
	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
)

func (c *client) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.Timestamp >= c.startTime {
		if msg, hasErrors := NewWhatsappBridgeMessage(c.wac, message).HasErrors(); !hasErrors {
			c.ExecuteSkill(msg, false)
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleImageMessage(message whatsapp.ImageMessage) {
	if message.Info.Timestamp >= c.startTime {
		if msg, hasErrors := NewWhatsappBridgeMessage(c.wac, message).HasErrors(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	if message.Info.Timestamp >= c.startTime {
		if msg, hasErrors := NewWhatsappBridgeMessage(c.wac, message).HasErrors(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleVideoMessage(message whatsapp.VideoMessage) {
	if message.Info.Timestamp >= c.startTime {
		if msg, hasErrors := NewWhatsappBridgeMessage(c.wac, message).HasErrors(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleAudioMessage(message whatsapp.AudioMessage) {
	if message.Info.Timestamp >= c.startTime {
		if msg, hasErrors := NewWhatsappBridgeMessage(c.wac, message).HasErrors(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}
