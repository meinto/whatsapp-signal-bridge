package whatsapp

import (
	"fmt"
	"os"
	"strings"

	"github.com/whatsapp-signal-bridge/bridge"
)

func (c *client) ExecuteSkill(msg bridge.Message, privateSkills bool) (executed bool, err error) {
	if strings.Contains(msg.Body(), "@bot") {
		if strings.Contains(msg.Body(), "help") {
			return c.ReplyHelp(msg)
		}

		if strings.Contains(msg.Body(), "status") {
			return c.ReplyStatus(msg)
		}

		if privateSkills {
			if strings.Contains(msg.Body(), "kill") && strings.Contains(msg.Body(), "whatsapp") {
				os.Exit(2)
				return true, nil
			}

			if strings.Contains(msg.Body(), "discover") {
				c.DiscoverChatIDHeader(msg)
			}
		}
	}
	return false, nil
}

func (c *client) ReplyHelp(msg bridge.Message) (executed bool, err error) {
	return c.Send(bridge.PlainTextMessage(strings.TrimSpace(`
available commands:

@bot status
> tells you the status of the bot
		`)).SetChatID(msg.ChatID()))
}

func (c *client) ReplyStatus(msg bridge.Message) (executed bool, err error) {
	if msg.ChatID() != "" {
		return c.Send(bridge.PlainTextMessage("i'm alive").SetChatID(msg.ChatID()))
	} else {
		return true, nil
	}
}

func (c *client) DiscoverChatIDHeader(msg bridge.Message) (executed bool, err error) {
	messageParts := strings.Split(msg.Body(), " ")
	last := messageParts[len(messageParts)-1]
	found := false
	contacts := c.wac.Store.Contacts
	for _, contact := range contacts {
		if strings.Contains(strings.ToLower(contact.Name), strings.ToLower(last)) || strings.Contains(contact.Jid, last) {
			found = true
			c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage(strings.TrimSpace(fmt.Sprintf(`
chatid:%s
chat:%s
---
						`, contact.Jid, contact.Name))))
		}
	}
	if !found {
		c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage(strings.TrimSpace(`nothing found :(`)))
	}
	return true, nil
}
