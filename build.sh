#!/bin/bash

REPO_NAME="Compel-Monitoring"
SERVER_NAME="compel-monitoring-server"
AGENT_NAME="compel-monitoring-agent"

echo "Building $SERVER_NAME"

if go build -o $GOPATH/bin/$SERVER_NAME -i github.com/adhuri/$REPO_NAME/$SERVER_NAME ;then
echo "+Successful"
else echo "-Failed"
fi


echo "Building $AGENT_NAME"

if go build -o $GOPATH/bin/$AGENT_NAME -i github.com/adhuri/$REPO_NAME/$AGENT_NAME ; then
echo "+Successful"
else echo "-Failed"
fi
