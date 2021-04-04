package message

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/whatsapp-signal-bridge/bridge"
)

type SignalBridgeMessage interface {
	bridge.Message
	Build() (message SignalBridgeMessage, err error)
}

type signalBridgeMessage struct {
	bridge.Message
	signalCLIMessage *SignalCLIMessage
}

func NewSignalBridgeMessage(signalCLIMessage *SignalCLIMessage) SignalBridgeMessage {
	signalBridgeMessage := &signalBridgeMessage{
		Message:          bridge.PlainMessage(),
		signalCLIMessage: signalCLIMessage,
	}
	return signalBridgeMessage
}

func (m *signalBridgeMessage) Build() (SignalBridgeMessage, error) {
	if m.signalCLIMessage.Envelope.DataMessage == nil {
		return nil, errors.New("Build: no data message")
	}

	m.Message.
		SetSender(m.signalCLIMessage.Envelope.Source).
		SetBody(m.signalCLIMessage.Envelope.DataMessage.Message)

	if m.signalCLIMessage.Envelope.DataMessage.Quote != nil {
		m.Message.
			SetQuote(&bridge.Quote{MessageType: bridge.TEXT_MESSAGE_TYPE, Body: &m.signalCLIMessage.Envelope.DataMessage.Quote.Text})
	}

	if len(m.signalCLIMessage.Envelope.DataMessage.Attachments) > 0 && m.signalCLIMessage.Envelope.DataMessage.Attachments[0].ID != "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			filePath := path.Join(
				homeDir,
				".local/share/signal-cli/attachments",
				m.signalCLIMessage.Envelope.DataMessage.Attachments[0].ID,
			)
			data, err := ioutil.ReadFile(filePath)
			defer func() {
				if err := os.Remove(filePath); err != nil {
					fmt.Println("error removing signal file", err)
				}
			}()
			if err == nil {
				m.Message.SetAttachment(&bridge.Attachment{
					Bytes: data,
					Type:  m.signalCLIMessage.Envelope.DataMessage.Attachments[0].ContentType,
				})
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	return m, nil
}
