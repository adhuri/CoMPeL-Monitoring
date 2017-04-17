#!/bin/bash

# Add influx key
curl -sL https://repos.influxdata.com/influxdb.key | sudo apt-key add -
source /etc/lsb-release
echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/influxdb.list


#Installation

sudo apt-get update && sudo apt-get install influxdb
echo "[INFO] Installation Success"

#Start
echo "[INFO]Starting influx db on default port"

sudo service influxdb start

#DB - square_holes

#influx

#CREATE DATABASE square_holes

echo "[NOTICE]MANUALLY CREATE DATABASE square_holes by 'CREATE DATABASE square_holes'"
