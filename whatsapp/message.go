package whatsapp

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"github.com/whatsapp-signal-bridge/bridge"
)

type WhatsappBridgeMessage interface {
	bridge.Message
	Build()
	SetInfo(whatsapp.MessageInfo) WhatsappBridgeMessage
	SetWhatsappQuote(*proto.Message) WhatsappBridgeMessage
	HasErrors() (WhatsappBridgeMessage, bool)
}

type whatsappBridgeMessage struct {
	bridge.Message
	wac    *whatsapp.Conn
	wam    interface{}
	errors []error
}

func NewWhatsappBridgeMessage(wac *whatsapp.Conn, whatsappMessage interface{}) WhatsappBridgeMessage {
	whatsappBridgeMessage := &whatsappBridgeMessage{
		Message: bridge.PlainMessage(),
		wac:     wac,
		wam:     whatsappMessage,
	}
	whatsappBridgeMessage.Build()
	return whatsappBridgeMessage
}

func (m *whatsappBridgeMessage) Build() {

	var attachmentBytes []byte
	var errAttachmentDownload error
	switch m.wam.(type) {
	case whatsapp.ImageMessage:
		wam := m.wam.(whatsapp.ImageMessage)
		attachmentBytes, errAttachmentDownload = wam.Download()
	case whatsapp.DocumentMessage:
		wam := m.wam.(whatsapp.DocumentMessage)
		attachmentBytes, errAttachmentDownload = wam.Download()
	case whatsapp.VideoMessage:
		wam := m.wam.(whatsapp.VideoMessage)
		attachmentBytes, errAttachmentDownload = wam.Download()
	case whatsapp.AudioMessage:
		wam := m.wam.(whatsapp.AudioMessage)
		attachmentBytes, errAttachmentDownload = wam.Download()
	}
	if errAttachmentDownload != nil {
		m.errors = append(m.errors, errAttachmentDownload)
	}

	if _, hasErrors := m.HasErrors(); !hasErrors {
		switch m.wam.(type) {
		case whatsapp.TextMessage:
			wam := m.wam.(whatsapp.TextMessage)
			m.SetInfo(wam.Info).
				SetWhatsappQuote(wam.ContextInfo.QuotedMessage).
				SetBody(wam.Text)

		case whatsapp.ImageMessage:
			wam := m.wam.(whatsapp.ImageMessage)
			m.SetInfo(wam.Info).
				SetWhatsappQuote(wam.ContextInfo.QuotedMessage).
				SetBody(wam.Caption).
				SetAttachment(&bridge.Attachment{
					Bytes: attachmentBytes,
					Type:  wam.Type,
				})

		case whatsapp.DocumentMessage:
			wam := m.wam.(whatsapp.DocumentMessage)
			m.SetInfo(wam.Info).
				SetWhatsappQuote(wam.ContextInfo.QuotedMessage).
				SetBody(wam.Title).
				SetAttachment(&bridge.Attachment{
					Bytes: attachmentBytes,
					Type:  wam.Type,
				})

		case whatsapp.VideoMessage:
			wam := m.wam.(whatsapp.VideoMessage)
			m.SetInfo(wam.Info).
				SetWhatsappQuote(wam.ContextInfo.QuotedMessage).
				SetBody(wam.Caption).
				SetAttachment(&bridge.Attachment{
					Bytes: attachmentBytes,
					Type:  wam.Type,
				})

		case whatsapp.AudioMessage:
			wam := m.wam.(whatsapp.AudioMessage)
			m.SetInfo(wam.Info).
				SetWhatsappQuote(wam.ContextInfo.QuotedMessage).
				SetAttachment(&bridge.Attachment{
					Bytes: attachmentBytes,
					Type:  wam.Type,
				})
		}
	}
}

func (m *whatsappBridgeMessage) SetInfo(info whatsapp.MessageInfo) WhatsappBridgeMessage {
	m.SetID(info.Id)
	m.SetChatID(info.RemoteJid)
	m.SetChatName(m.getChatName(info))
	m.SetSender(m.getSender(info))
	return m
}

func (m *whatsappBridgeMessage) SetWhatsappQuote(quotedMessage *proto.Message) WhatsappBridgeMessage {
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

func (m *whatsappBridgeMessage) getChatName(info whatsapp.MessageInfo) string {
	chatID := m.ChatID()
	return m.wac.Store.Chats[chatID].Name
}

func (m *whatsappBridgeMessage) getSender(info whatsapp.MessageInfo) string {
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

func (m *whatsappBridgeMessage) getSenderID(info whatsapp.MessageInfo) (senderID string) {
	if info.Source.Participant == nil {
		senderID = info.RemoteJid
	} else {
		senderID = *info.Source.Participant
	}
	return senderID
}

func (m *whatsappBridgeMessage) HasErrors() (WhatsappBridgeMessage, bool) {
	return m, len(m.errors) > 0
}
