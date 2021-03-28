#!/bin/bash

if [[ $1 == "" ]]; then
  echo "signal bot number missing"
  echo "start with: ./start.sh <bot-number> <receiver-number>"
  exit
fi

if [[ $2 == "" ]]; then
  echo "receiver number missing"
  echo "start with: ./start.sh <bot-number> <receiver-number>"
  exit
fi

go build -o bot .
pkill -f whatsappSignalBridge
nohup bash -c "exec -a whatsappSignalBridge $(pwd)/bot --bot=$1 --receiver=$2" > /dev/null 2>&1 & 