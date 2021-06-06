package whatsapp

import (
	"errors"
	"os"

	"github.com/whatsapp-signal-bridge/bridge"
)

var ErrRestoreSessionConnectionTimeout = errors.New("restore session connection timed out")

func (c *client) HandleError(err error) {
	c.Publish(bridge.WHATSAPP_QUEUE, bridge.ErrorMessage(err))

	// Service will restart automatically
	// see ./service/
	os.Exit(2)
}
