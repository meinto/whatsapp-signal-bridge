package whatsapp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
)

type WhatsappClient interface {
	bridge.Client
}

type client struct {
	bridge.Client
	wac            *whatsapp.Conn
	restoreAttemts int
	startTime      uint64
}

type WhatsappClientOptions struct {
	Queue bridge.Queue
}

func StartClient(options WhatsappClientOptions) {
	wac, _ := whatsappLogin()

	c := &client{
		bridge.NewClient(options.Queue),
		wac,
		0,
		uint64(time.Now().Unix()),
	}

	c.Subscribe(bridge.SIGNAL_QUEUE, func(msg bridge.Message) {
		if executed, err := c.ExecuteSkill(msg, true); !executed || err != nil {
			c.Replay(msg)
		}
	})

	c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainMessage().SetBody("service started"))

	wac.AddHandler(c)
}

func (c *client) Send(msg bridge.Message) (executed bool, err error) {
	waMsg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: msg.ChatID(),
		},
		Text: msg.Body(),
	}
	if _, err := c.wac.Send(waMsg); err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
		return false, err
	}
	return true, nil
}

func (c *client) Replay(msg bridge.Message) (executed bool, err error) {
	if msg.Quote() != nil {
		quoteBody := msg.Quote().Body
		if quoteBody == nil {
			return
		}

		quoteParts := strings.Split(*quoteBody, "---")
		header := quoteParts[0]

		type whatsappMetaData struct {
			// WhatsappMessageID string
			WhatsappChatID string
		}
		metaData := &whatsappMetaData{}
		headerParts := strings.Split(header, "\n")
		for _, row := range headerParts {
			rowParts := strings.Split(row, ":")
			if len(rowParts) > 1 {
				key := rowParts[0]
				val := rowParts[1]

				switch strings.TrimSpace(key) {
				// case "id":
				// 	metaData.WhatsappMessageID = strings.TrimSpace(val)
				case "chatid":
					metaData.WhatsappChatID = strings.TrimSpace(val)
				}
			}
		}

		msgText := msg.Body()
		if len(quoteParts) > 1 {
			body := strings.TrimSpace(quoteParts[1])
			bodyParts := strings.Split(body, "\n")
			quoteTextParts := []string{}
			for _, p := range bodyParts {
				quoteTextParts = append(quoteTextParts, "> "+p)
			}
			quoteText := strings.Join(quoteTextParts, "\n")
			if strings.Contains(strings.ToLower(msgText), "quote") {
				msgText = strings.Replace(msgText, "quote", quoteText, -1)
				msgText = strings.Replace(msgText, "Quote", quoteText, -1)
			}
		}

		waMessageInfo := whatsapp.MessageInfo{
			RemoteJid: metaData.WhatsappChatID,
		}

		var waMsg interface{}
		if msg.Attachment() == nil {
			waMsg = whatsapp.TextMessage{
				Info: waMessageInfo,
				Text: msgText,
			}
		}

		if msg.Attachment() != nil {
			if strings.Contains(msg.Attachment().Type, "image") {
				waMsg = whatsapp.ImageMessage{
					Info:    waMessageInfo,
					Type:    msg.Attachment().Type,
					Content: bytes.NewReader(msg.Attachment().Bytes),
					Caption: msgText,
				}
			}

			if strings.Contains(msg.Attachment().Type, "video") {
				waMsg = whatsapp.VideoMessage{
					Info:    waMessageInfo,
					Type:    msg.Attachment().Type,
					Content: bytes.NewReader(msg.Attachment().Bytes),
					Caption: msgText,
				}
			}

			if strings.Contains(msg.Attachment().Type, "document") {
				waMsg = whatsapp.DocumentMessage{
					Info:    waMessageInfo,
					Type:    msg.Attachment().Type,
					Content: bytes.NewReader(msg.Attachment().Bytes),
					Title:   msgText,
				}
			}

			if strings.Contains(msg.Attachment().Type, "audio") {
				waMsg = whatsapp.AudioMessage{
					Info:    waMessageInfo,
					Type:    msg.Attachment().Type,
					Content: bytes.NewReader(msg.Attachment().Bytes),
				}
			}
		}

		if _, err := c.wac.Send(waMsg); err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (c *client) HandleError(err error) {
	c.Publish(bridge.WHATSAPP_QUEUE, bridge.ErrorMessage(err))

	// https://github.com/Rhymen/go-whatsapp/issues/343#issuecomment-621584947
	if strings.Contains(err.Error(), "server closed connection") || strings.Contains(err.Error(), "close 1006") {
		c.restoreAttemts++
		c.startTime = uint64(time.Now().Unix())

		if c.restoreAttemts > 5 {
			c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage("Cannot restore session.\nPlease restart manually."))
		}

		c.RestoreWhatsappConnection()
	}
}

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
