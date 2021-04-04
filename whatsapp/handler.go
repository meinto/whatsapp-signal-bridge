package whatsapp

import (
	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
	"github.com/whatsapp-signal-bridge/whatsapp/message"
)

func (c *client) HandleTextMessage(wam whatsapp.TextMessage) {
	if wam.Info.Timestamp >= c.startTime {
		if msg, hasErrors := message.NewWhatsappBridgeMessage(c.wac, wam).Build(); !hasErrors {
			c.ExecuteSkill(msg, false)
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleImageMessage(wam whatsapp.ImageMessage) {
	if wam.Info.Timestamp >= c.startTime {
		if msg, hasErrors := message.NewWhatsappBridgeMessage(c.wac, wam).Build(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleDocumentMessage(wam whatsapp.DocumentMessage) {
	if wam.Info.Timestamp >= c.startTime {
		if msg, hasErrors := message.NewWhatsappBridgeMessage(c.wac, wam).Build(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleVideoMessage(wam whatsapp.VideoMessage) {
	if wam.Info.Timestamp >= c.startTime {
		if msg, hasErrors := message.NewWhatsappBridgeMessage(c.wac, wam).Build(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}

func (c *client) HandleAudioMessage(wam whatsapp.AudioMessage) {
	if wam.Info.Timestamp >= c.startTime {
		if msg, hasErrors := message.NewWhatsappBridgeMessage(c.wac, wam).Build(); !hasErrors {
			c.Publish(bridge.WHATSAPP_QUEUE, msg)
		}
	}
}
