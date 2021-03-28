package whatsapp

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"github.com/whatsapp-signal-bridge/bridge"
)

type WhatsappMessage interface {
	bridge.Message
	SetInfo(whatsapp.MessageInfo) WhatsappMessage
	SetWhatsappQuote(*proto.Message) WhatsappMessage
}

type whatsappMessage struct {
	bridge.Message
	wac *whatsapp.Conn
}

func NewWhatsappMessage(wac *whatsapp.Conn) WhatsappMessage {
	return &whatsappMessage{
		Message: bridge.PlainMessage(),
		wac:     wac,
	}
}

func (m *whatsappMessage) SetInfo(info whatsapp.MessageInfo) WhatsappMessage {
	m.SetID(info.Id)
	m.SetChatID(info.RemoteJid)
	m.SetChatName(m.getChatName(info))
	m.SetSender(m.getSender(info))
	return m
}

func (m *whatsappMessage) SetWhatsappQuote(quotedMessage *proto.Message) WhatsappMessage {
	if quotedMessage != nil {
		if quotedMessage.Conversation != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.TEXT_MESSAGE_TYPE,
				Body:        quotedMessage.Conversation,
			})
		}

		if quotedMessage.ImageMessage != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.IMAGE_MESSAGE_TYPE,
				Body:        quotedMessage.ImageMessage.Caption,
			})
		}

		if quotedMessage.DocumentMessage != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.DOCUMENT_MESSAGE_TYPE,
				Body:        quotedMessage.DocumentMessage.Title,
			})
		}

		if quotedMessage.DocumentMessage != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.DOCUMENT_MESSAGE_TYPE,
				Body:        quotedMessage.DocumentMessage.Title,
			})
		}

		if quotedMessage.VideoMessage != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.VIDEO_MESSAGE_TYPE,
				Body:        quotedMessage.VideoMessage.Caption,
			})
		}

		if quotedMessage.AudioMessage != nil {
			m.SetQuote(&bridge.Quote{
				MessageType: bridge.AUDIO_MESSAGE_TYPE,
			})
		}
	}
	return m
}

func (m *whatsappMessage) getChatName(info whatsapp.MessageInfo) string {
	chatID := m.ChatID()
	return m.wac.Store.Chats[chatID].Name
}

func (m *whatsappMessage) getSender(info whatsapp.MessageInfo) string {
	if info.FromMe {
		return "Me"
	}

	senderID := m.getSenderID(info)
	sender := m.wac.Store.Contacts[senderID].Name
	if sender == "" {
		if jidParts := strings.Split(m.wac.Store.Contacts[senderID].Jid, "@"); len(jidParts) > 0 {
			sender = jidParts[0]
		}
	}

	return sender
}

func (m *whatsappMessage) getSenderID(info whatsapp.MessageInfo) (senderID string) {
	if info.Source.Participant == nil {
		senderID = info.RemoteJid
	} else {
		senderID = *info.Source.Participant
	}
	return senderID
}
