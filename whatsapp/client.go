package whatsapp

import (
	"fmt"
	"os"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/whatsapp-signal-bridge/bridge"
	"github.com/whatsapp-signal-bridge/logger"
	"github.com/whatsapp-signal-bridge/whatsapp/message"
)

type WhatsappClient interface {
	bridge.Client
}

type client struct {
	bridge.Client
	wac            *whatsapp.Conn
	restoreAttemts int
	startTime      uint64
	logger         logger.Logger
}

type WhatsappClientOptions struct {
	Queue bridge.Queue
}

func StartClient(options WhatsappClientOptions) {
	wac, _ := whatsappLogin()

	c := &client{
		bridge.NewClient(options.Queue, "whatsapp"),
		wac,
		0,
		uint64(time.Now().Unix()),
		logger.NewLogger("whatsapp", logger.LOG_LEVEL_DEBUG),
	}

	c.Subscribe(bridge.SIGNAL_QUEUE, func(msg bridge.Message) {
		if executed, err := c.ExecuteSkill(msg, true); !executed || err != nil {
			c.Send(msg)
		}
	})

	c.Publish(bridge.WHATSAPP_QUEUE, bridge.PlainMessage().SetBody("service started"))

	wac.AddHandler(c)
}

func (c *client) Send(msg bridge.Message) (executed bool, err error) {
	if wam, err := message.NewWhatsappMessage(msg).Build(); err == nil {
		if _, err := c.wac.Send(wam); err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
			return false, err
		}
	}
	return false, nil
}
