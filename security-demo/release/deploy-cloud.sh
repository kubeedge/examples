#!/bin/sh 

COMMAND_PATH=/opt/spire

clean()
{
  pkill spire-server; pkill spire-agent; pkill spiffe-helper; pkill ghostunnel; pkill cloud-app
  rm .data/*
  rm certs/*.p12 certs/*.pem
}

check_command_status()
{
  if [ $2 -ne 0 ]; then
    echo $return
#TBD define exit codes for command failures
    exit 1
  else 
    echo "command $1 is successful\n"
  fi
}

case $1 in 
  "clean")
     clean
     exit 0
     ;;
  *)
     ;;
esac

$COMMAND_PATH/commands.sh start-cloud-server
check_command_status start-server $?

echo "Waiting for server to be initialized ..."
sleep 5

$COMMAND_PATH/commands.sh cloud-node
check_command_status cloud-node $?

echo "Waiting for node to be initialized ..."
sleep 5

$COMMAND_PATH/commands.sh cloud-app cloud-hub unix:uid:1000
check_command_status cloud-app $?

$COMMAND_PATH/commands.sh start-cloud-hub
check_command_status start-cloud-hub $?

echo "Waiting for hub initialization ..."
sleep 5

#TBD : this is a hack for keystore requirement
pkill spiffe-helper; pkill ghostunnel; pkill cloud-app

$COMMAND_PATH/commands.sh start-cloud-hub
check_command_status start-cloud-hub $?

