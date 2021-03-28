package signal

import (
	"strings"

	"github.com/whatsapp-signal-bridge/bridge"
)

func (c *client) ExecuteSkill(msg bridge.Message) (executed bool, err error) {
	if strings.Contains(msg.Body(), "@bot") {
		if strings.Contains(msg.Body(), "help") {
			return c.ReplyHelp()
		}
	}
	return false, nil
}

func (c *client) ReplyHelp() (executed bool, err error) {
	return c.Send(bridge.PlainTextMessage(strings.TrimSpace(`
available commands:

@bot restore whatsapp
> restores whatsapp connection

@bot discover <substring>
> list of possible chats which match the substring

---
Info: 

The Keyword "quote" will be replaced with the quoted message text
		`)))
}
