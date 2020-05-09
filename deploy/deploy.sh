#!/bin/bash

eval "$(cat ./deploy/.env <(echo) <(declare -x))"

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
scp -P ${PORT} twchatbot config.yml deploy/twchatbot.conf deploy/twchatbot.service ${USER_NAME}@${HOST}:/home/${USER_NAME}/
ssh ${USER_NAME}@${HOST} -p ${PORT} << EOS
sudo mv /home/${USER_NAME}/twchatbot.conf /etc/nginx/conf.d/
sudo mv /home/${USER_NAME}/twchatbot.service /etc/systemd/system/
sudo systemctl enable twchatbot.service
sudo systemctl restart twchatbot.service
sudo nginx -s reload
EOS
