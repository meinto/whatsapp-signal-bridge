package message

import (
	"bytes"
	"errors"
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
)

type WhatsappMessage interface {
	Build() (interface{}, error)
	GetMetaData() (*metaData, error)
	GetMessageBody() (string, error)
}

type metaData struct {
	// WhatsappMessageID string
	WhatsappChatID *string
}

type whatsappMessage struct {
	bridgeMessage bridge.Message
}

func NewWhatsappMessage(bridgeMessage bridge.Message) WhatsappMessage {
	return &whatsappMessage{bridgeMessage: bridgeMessage}
}

func (m *whatsappMessage) Build() (wam interface{}, err error) {
	md, err := m.GetMetaData()
	if err != nil {
		return nil, err
	}

	body, err := m.GetMessageBody()
	if err != nil {
		return nil, err
	}

	messageInfo := whatsapp.MessageInfo{
		RemoteJid: *md.WhatsappChatID,
	}

	if m.bridgeMessage.Attachment() == nil {
		wam = whatsapp.TextMessage{
			Info: messageInfo,
			Text: body,
		}
	}

	if m.bridgeMessage.Attachment() != nil {
		if strings.Contains(m.bridgeMessage.Attachment().Type, "image") {
			wam = whatsapp.ImageMessage{
				Info:    messageInfo,
				Type:    m.bridgeMessage.Attachment().Type,
				Content: bytes.NewReader(m.bridgeMessage.Attachment().Bytes),
				Caption: body,
			}
		}

		if strings.Contains(m.bridgeMessage.Attachment().Type, "video") {
			wam = whatsapp.VideoMessage{
				Info:    messageInfo,
				Type:    m.bridgeMessage.Attachment().Type,
				Content: bytes.NewReader(m.bridgeMessage.Attachment().Bytes),
				Caption: body,
			}
		}

		if strings.Contains(m.bridgeMessage.Attachment().Type, "document") {
			wam = whatsapp.DocumentMessage{
				Info:    messageInfo,
				Type:    m.bridgeMessage.Attachment().Type,
				Content: bytes.NewReader(m.bridgeMessage.Attachment().Bytes),
				Title:   body,
			}
		}

		if strings.Contains(m.bridgeMessage.Attachment().Type, "audio") {
			wam = whatsapp.AudioMessage{
				Info:    messageInfo,
				Type:    m.bridgeMessage.Attachment().Type,
				Content: bytes.NewReader(m.bridgeMessage.Attachment().Bytes),
			}
		}
	}

	return wam, nil
}

func (m *whatsappMessage) GetMetaData() (*metaData, error) {
	if m.bridgeMessage.Quote() == nil {
		return nil, errors.New("BuildMetaData: quote missing")
	}

	quoteBody := m.bridgeMessage.Quote().Body
	if quoteBody == nil {
		return nil, errors.New("BuildMetaData: quote body missing")
	}

	quoteParts := strings.Split(*quoteBody, "---")
	header := quoteParts[0]

	md := &metaData{}
	headerParts := strings.Split(header, "\n")

	for _, row := range headerParts {
		rowParts := strings.Split(row, ":")
		if len(rowParts) > 1 {
			key := rowParts[0]
			val := strings.TrimSpace(rowParts[1])

			switch strings.TrimSpace(key) {
			// case "id":
			// 	md.WhatsappMessageID = strings.TrimSpace(val)
			case "chatid":
				md.WhatsappChatID = &val
			}
		}
	}

	if md.WhatsappChatID == nil {
		return nil, errors.New("BuildMetaData: WhatsappChatID missing")
	}

	return md, nil
}

func (m *whatsappMessage) GetMessageBody() (string, error) {
	if m.bridgeMessage.Quote() == nil {
		return "", errors.New("GetMessageBody: quote missing")
	}

	quoteBody := m.bridgeMessage.Quote().Body
	if quoteBody == nil {
		return "", errors.New("BuildMetaData: quote body missing")
	}

	quoteParts := strings.Split(*quoteBody, "---")

	msgText := m.bridgeMessage.Body()
	if len(quoteParts) > 1 {
		quoteBody := strings.TrimSpace(quoteParts[1])
		quoteBodyParts := strings.Split(quoteBody, "\n")
		quoteTextParts := []string{}
		for _, p := range quoteBodyParts {
			if !strings.HasPrefix(p, "▒") {
				quoteTextParts = append(quoteTextParts, "▒ _"+p+"_")
			}
		}
		quoteText := strings.Join(quoteTextParts, "\n")
		if strings.Contains(strings.ToLower(msgText), "quote") {
			msgText = strings.Replace(msgText, "quote", quoteText, -1)
			msgText = strings.Replace(msgText, "Quote", quoteText, -1)
		}
	}

	return msgText, nil
}
