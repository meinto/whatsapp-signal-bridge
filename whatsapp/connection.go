package whatsapp

import "github.com/whatsapp-signal-bridge/bridge"

func (c *client) RestoreWhatsappConnection() error {
	c.wac.Disconnect()
	if err := c.wac.Restore(); err != nil {
		c.Publish(bridge.WHATSAPP_QUEUE, bridge.ErrorMessage(err, "error restoring whatsapp session"))
		return err
	} else {
		c.restoreAttemts = 0
		c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage("whatsapp connection restored"))
		return nil
	}
}
