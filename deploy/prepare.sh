#!/bin/bash

eval "$(cat ./deploy/.env <(echo) <(declare -x))"

scp -P ${PORT} deploy/twchatbot.conf deploy/twchatbot.service ${USER_NAME}@${HOST}:/home/${USER_NAME}/
ssh -t ${USER_NAME}@${HOST} -p ${PORT} << EOS
sudo mv /home/${USER_NAME}/twchatbot.conf /etc/nginx/conf.d/
sudo mv /home/${USER_NAME}/twchatbot.service /etc/systemd/system/
sudo systemctl enable twchatbot.service
/usr/local/certbot/certbot-auto certonly --standalone -d ${HOST} -m ${MAIL} --agree-tos -n
EOS
