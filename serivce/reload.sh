#!/bin/bash

cd ..
go build -o bot .
sudo service whatsapp-signal-bridge stop
sudo systemctl disable whatsapp-signal-bridge.service
sudo systemctl enable whatsapp-signal-bridge.service
sudo service whatsapp-signal-bridge start