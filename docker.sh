#!/bin/bash
cd /app

if [ -z "$DOCKER_DEBUG" ]
then
  npm start
else
  echo 'Starting bot with inspector exposed at 0.0.0.0:9229'
  npm run debug
fi