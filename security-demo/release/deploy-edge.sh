#!/bin/sh -x

COMMAND_PATH=/opt/spire

clean()
{
	pkill spire-server; pkill spire-agent; pkill spiffe-helper; pkill ghostunnel; sudo pkill edge_core;
	rm /opt/spire/.dserverdata/*
# agent crashes with old certificates ??
	rm /opt/spire/.data/*
	rm /opt/spire/certs/*.pem /opt/spire/certs/*.p12
	rm /opt/spire/user-app/certs/*.pem /opt/spire/user-app/certs/*.p12
	rm /opt/spire/event-bus/certs/*.pem /opt/spire/event-bus/certs/*.p12
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
  sleep 5
}

case $1 in 
  clean) 
    clean
    exit 0
    ;;
esac

# register edge node to upstream server
$COMMAND_PATH/commands.sh edge-node
check_command_status edge-node $?


#create entry for edge-hub 
$COMMAND_PATH/commands.sh edge-control-app edge-hub unix:uid:1000
check_command_status edge-app $?


#start downstream edge server- uses edge-hub certificates 
$COMMAND_PATH/commands.sh start-edge-server
check_command_status start-edge-server $?


# register edge app node to edge server
$COMMAND_PATH/commands.sh edge-app-node
check_command_status edge-app-node $?


#create entry for event-bus
$COMMAND_PATH/commands.sh edge-app event-bus unix:uid:1000
check_command_status event-bus $?


#create entry for ligthmapper; this will require identifying unqiue process, run as different user id ? 
$COMMAND_PATH/commands.sh edge-app lightmapper unix:uid:1000
check_command_status lightmapper $?


#start edge-hub 
$COMMAND_PATH/commands.sh start-edge-hub
check_command_status start-edge-hub $?


#start edge event-bus - helper and tunnel only
$COMMAND_PATH/commands.sh start-edge-event-bus
check_command_status start-edge-event-bus $?


#start edge light mapper app
$COMMAND_PATH/commands.sh start-edge-app
check_command_status start-edge-app $?


#TBD : this is a hack for keystore requirement
pkill spiffe-helper; pkill ghostunnel; pkill lightmapper; sudo pkill edge_core;
sleep 5

$COMMAND_PATH/commands.sh start-edge-hub
check_command_status start-edge-hub $?


$COMMAND_PATH/commands.sh start-edge-event-bus
check_command_status start-event-bus $?


$COMMAND_PATH/commands.sh start-edge-app
check_command_status start-edge-app $?

