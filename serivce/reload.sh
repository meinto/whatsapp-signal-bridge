#!/bin/bash

go build -o bot .
sudo service whatsapp-signal-bridge start
sudo systemctl disable whatsapp-signal-bridge.service
sudo systemctl enable whatsapp-signal-bridge.service
sudo service whatsapp-signal-bridge start