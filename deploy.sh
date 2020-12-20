#!/bin/bash

cd infra || exit
npm i
cdk deploy
now=$(date +"%x - %T")
echo "Last run : $now"