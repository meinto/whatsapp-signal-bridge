package whatsapp

import (
	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
)

func (c *client) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.Timestamp >= c.startTime {
		msg := NewWhatsappMessage(c.wac).
			SetInfo(message.Info).
			SetWhatsappQuote(message.ContextInfo.QuotedMessage).
			SetBody(message.Text)

		c.ExecuteSkill(msg, false)
		c.Publish(bridge.WHATSAPP_QUEUE, msg)
	}
}

func (c *client) HandleImageMessage(message whatsapp.ImageMessage) {
	if message.Info.Timestamp >= c.startTime {
		if bytes, err := message.Download(); err == nil {
			c.Publish(bridge.WHATSAPP_QUEUE, NewWhatsappMessage(c.wac).
				SetInfo(message.Info).
				SetWhatsappQuote(message.ContextInfo.QuotedMessage).
				SetBody(message.Caption).
				SetAttachment(&bridge.Attachment{
					Bytes: bytes,
					Type:  message.Type,
				}),
			)
		}
	}
}

func (c *client) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	if message.Info.Timestamp >= c.startTime {
		if bytes, err := message.Download(); err == nil {
			c.Publish(bridge.WHATSAPP_QUEUE, NewWhatsappMessage(c.wac).
				SetInfo(message.Info).
				SetWhatsappQuote(message.ContextInfo.QuotedMessage).
				SetBody(message.Title).
				SetAttachment(&bridge.Attachment{
					Bytes: bytes,
					Type:  message.Type,
				}),
			)
		}
	}
}

func (c *client) HandleVideoMessage(message whatsapp.VideoMessage) {
	if message.Info.Timestamp >= c.startTime {
		if bytes, err := message.Download(); err == nil {
			c.Publish(bridge.WHATSAPP_QUEUE, NewWhatsappMessage(c.wac).
				SetInfo(message.Info).
				SetWhatsappQuote(message.ContextInfo.QuotedMessage).
				SetBody(message.Caption).
				SetAttachment(&bridge.Attachment{
					Bytes: bytes,
					Type:  message.Type,
				}),
			)
		}
	}
}

func (c *client) HandleAudioMessage(message whatsapp.AudioMessage) {
	if message.Info.Timestamp >= c.startTime {
		if bytes, err := message.Download(); err == nil {
			c.Publish(bridge.WHATSAPP_QUEUE, NewWhatsappMessage(c.wac).
				SetInfo(message.Info).
				SetWhatsappQuote(message.ContextInfo.QuotedMessage).
				SetAttachment(&bridge.Attachment{
					Bytes: bytes,
					Type:  message.Type,
				}),
			)
		}
	}
}
