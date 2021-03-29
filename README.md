# Whatsapp Signal Bridge

This application lets you forward all messages which you receive in your Whatsapp account to Signal. Big thanks to all contributors of the [Rhymen/go-whatsapp](https://github.com/Rhymen/go-whatsapp) and [signal-cli](https://github.com/AsamK/signal-cli).

## How it works

For this application to work you need a separate phone number for your bot. The bot sends the Whatsapp messages to your main Signal number. So you receive all Whatsapp messages in a separate Signal chat.

The messages you receive in this Signal chat include the Whatsapp chat-id, the name or the number of the sender, the group name, the message and the attachments.

_Example text message structure:_

```
chatid:000000-000000@g.us
chat:Chat with friends
sender:Mike
---
Hi there
```

You can reply to the messages by quoting them. The reply will be sent back through the bot in the corresponding Whatsapp chat including attachments.

## Requirements

You need to install [signal-cli](https://github.com/AsamK/signal-cli).

([Native Library](https://github.com/AsamK/signal-cli/wiki/Provide-native-lib-for-libsignal)) ([Use Correct Libsignal-client Version](https://github.com/AsamK/signal-cli/issues/562#issuecomment-792274797)) ([Maybe useful](https://github.com/AsamK/signal-cli/discussions/393#discussioncomment-246169))

After installation you have to [register your separate phone number](https://github.com/AsamK/signal-cli/blob/master/man/signal-cli.1.adoc#register).

## Build

```go
go build -o bot .
```

## First Run

On the first run you have to connect your Whatsapp account to the application. This is done by adding the application as Whatsapp-Web client.

On the first run you have to scan the qr-code which will be presented on your terminal.

```bash
./bot --bot=<bot-number> --receiver=<signal-receiver-number>
# now scan the presented qr-code
```

After the qr-code was scanned the connection should be established and the session is stored in a file. On the next run of this application, the session will be loaded directly. You don't have to scan the qr-code again.

## Start Bridge for Production

```bash
# this script starts the application in the background
./start.sh <bot-numer> <signal-receiver-number>
# stopping the background script
./stop.sh
```

> Remember that only one concurrent Whatsapp-Web connection can be active!

## Skills

You can control the application with a couple of commands. Type `@bot help` in Signal or Whatsapp for more information.

## Legal

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Whatsapp or any of its affiliates or subsidiaries. This is an independent and unofficial software. Use at your own risk.

## License

The MIT License (MIT)

Copyright (c) 2021

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
