package signal

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"os/exec"
	"path"

	"github.com/whatsapp-signal-bridge/bridge"
)

type SignalClient interface {
	bridge.Client
	receiveMessages()
}

type client struct {
	bridge.Client
	botNumber      string
	receiverNumber string
}

type SignalClientOptions struct {
	Queue          bridge.Queue
	BotNumber      string
	ReceiverNumber string
}

func StartClient(options SignalClientOptions) {
	c := &client{
		bridge.NewClient(options.Queue),
		options.BotNumber,
		options.ReceiverNumber,
	}
	c.Subscribe(bridge.WHATSAPP_QUEUE, func(msg bridge.Message) {
		c.Send(msg)
	})
	go c.receiveMessages()
}

func (c *client) Send(msg bridge.Message) (executed bool, err error) {
	textMessage := ""
	// if msg.ID() != "" {
	// 	textMessage += "id:" + msg.ID() + "\n"
	// }
	if msg.ChatID() != "" {
		textMessage += "chatid:" + msg.ChatID() + "\n"
	}
	if msg.ChatName() != "" {
		textMessage += "chat:" + msg.ChatName() + "\n"
	}
	if msg.Sender() != "" {
		textMessage += "sender:" + msg.Sender() + "\n"
	}
	if textMessage != "" {
		textMessage += "---\n"
	}

	if msg.Quote() != nil {
		textMessage += "> " + string(msg.Quote().MessageType)
		if quoteText := msg.Quote().Body; quoteText != nil {
			textMessage += ":" + *quoteText
		}
		textMessage += "\n"
	}
	textMessage += msg.Body()

	cmd := exec.Command("signal-cli", "-u", c.botNumber, "send", c.receiverNumber, "-m", textMessage)

	if msg.Attachment() != nil && msg.Attachment().Bytes != nil {
		extensions, err := mime.ExtensionsByType(msg.Attachment().Type)
		if extensions != nil && err == nil {
			filePath := path.Join(os.TempDir(), msg.ID()+extensions[0])
			ioutil.WriteFile(filePath, msg.Attachment().Bytes, 0755)
			cmd.Args = append(cmd.Args, "-a", filePath)
			defer func() {
				if err := os.Remove(filePath); err != nil {
					log.Println("error removing file:", filePath)
				}
			}()
		}
	}

	if err := cmd.Start(); err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		return false, err
	}
	return true, nil
}

func (c *client) receiveMessages() {
	cmd := exec.Command("signal-cli", "-u", c.botNumber, "receive", "--json")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println(err)
	}
	cmd.Start()

	currentMessage := bridge.PlainMessage()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		row := scanner.Text()

		log.Println(row)

		var signalCLIMessage SignalCLIMessage
		err := json.Unmarshal([]byte(row), &signalCLIMessage)
		if err != nil {
			log.Println(err)
		}

		if signalCLIMessage.Envelope.DataMessage != nil {
			currentMessage.
				SetSender(signalCLIMessage.Envelope.Source).
				SetBody(signalCLIMessage.Envelope.DataMessage.Message)

			if signalCLIMessage.Envelope.DataMessage.Quote != nil {
				currentMessage.
					SetQuote(&bridge.Quote{MessageType: bridge.TEXT_MESSAGE_TYPE, Body: &signalCLIMessage.Envelope.DataMessage.Quote.Text})
			}

			if len(signalCLIMessage.Envelope.DataMessage.Attachments) > 0 && signalCLIMessage.Envelope.DataMessage.Attachments[0].ID != "" {
				homeDir, err := os.UserHomeDir()
				if err == nil {
					filePath := path.Join(
						homeDir,
						".local/share/signal-cli/attachments",
						signalCLIMessage.Envelope.DataMessage.Attachments[0].ID,
					)
					data, err := ioutil.ReadFile(filePath)
					defer func() {
						if err := os.Remove(filePath); err != nil {
							log.Println("error removing signal file")
						}
					}()
					if err == nil {
						currentMessage.SetAttachment(&bridge.Attachment{
							Bytes: data,
							Type:  signalCLIMessage.Envelope.DataMessage.Attachments[0].ContentType,
						})
					} else {
						log.Println(err)
					}
				} else {
					log.Println(err)
				}
			}

			if executed, err := c.ExecuteSkill(currentMessage); !executed || err != nil {
				c.Publish(bridge.SIGNAL_QUEUE, currentMessage)
			}
		}
	}

	cmd.Wait()
	go c.receiveMessages()
}
