package whatsapp

import (
	"errors"
	"strings"
	"time"

	"github.com/whatsapp-signal-bridge/bridge"
)

var ErrRestoreSessionConnectionTimeout = errors.New("restore session connection timed out")

func (c *client) HandleError(err error) {
	c.Publish(bridge.WHATSAPP_QUEUE, bridge.ErrorMessage(err))

	// https://github.com/Rhymen/go-whatsapp/issues/343#issuecomment-621584947
	if strings.Contains(err.Error(), "server closed connection") || strings.Contains(err.Error(), "close 1006") {
		c.startTime = uint64(time.Now().Unix())
		c.RestoreWhatsappConnection()
	}
}

func (c *client) RestoreWhatsappConnection() error {
	c.wac.Disconnect()
	if err := c.wac.Restore(); err != nil {
		if err.Error() == ErrRestoreSessionConnectionTimeout.Error() {
			c.Publish(bridge.WHATSAPP_QUEUE, bridge.ErrorMessage(err, "error restoring whatsapp session"))
			c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage("please restart your whatsapp on your smartphone"))
		}
		return err
	} else {
		c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainTextMessage("whatsapp connection restored"))
		return nil
	}
}
