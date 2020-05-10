#!/bin/bash

eval "$(cat ./deploy/.env <(echo) <(declare -x))"

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
ssh -t ${USER_NAME}@${HOST} -p ${PORT} "sudo systemctl stop twchatbot.service"
scp -P ${PORT} twchatbot config.yml ${USER_NAME}@${HOST}:/home/${USER_NAME}/
ssh -t ${USER_NAME}@${HOST} -p ${PORT} "sudo systemctl restart twchatbot.service"
