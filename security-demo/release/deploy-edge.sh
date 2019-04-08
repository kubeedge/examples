#!/bin/bash

COMMAND_PATH=/opt/spire

clean()
{
	pkill spire-server; pkill spire-agent; pkill spiffe-helper; pkill ghostunnel; sudo pkill edge_core;
# agent crashes with old certificates ??
	rm /opt/spire/.data/*
  rm /opt/spire/.agentdata/*
  rm /opt/spire/.app-data/*
	rm /opt/spire/certs/*
	rm /opt/spire/user-app/certs/*
	rm /opt/spire/event-bus/certs/*
}

check_command_status()
{
  if [ $2 -ne 0 ]; then
    echo $return
#TBD define exit codes for command failures
    exit 1
  else 
    echo "command $1 is successful" >> deploy-edge.log
  fi
  sleep 5
}

case $1 in 
  clean) 
    clean
    exit 0
    ;;
esac

date=`date`
echo "**Started at $date" | tee -a deploy-edge.log
echo ""

echo "*Generating configurations" | tee -a deploy-edge.log
echo ""
$COMMAND_PATH/generate-edge-config.sh

echo "*Start SPIRE agent at edge node which registers and attests edge node to cloud " | tee -a deploy-edge.log
echo ""
# register edge node to upstream server
$COMMAND_PATH/commands.sh edge-node &>> deploy-edge.log
check_command_status edge-node $?
echo "..Waiting for edge node to be initialized .." | tee -a deploy-edge.log
echo ""
sleep 5

echo "*Start edge core components. Components include sidecars (Spiffe-helper - for certificate acquisition rotation, ghostunnel - for mTLS tunnel) " | tee -a deploy-edge.log
echo "*Workload attestation is performed against the unique selectors \"unix:uid:1000\"" | tee -a deploy-edge.log
echo "*All processes running under uid 1000 are attested and provided leaf certificates for communication with the help of spiffe-helper" | tee -a deploy-edge.log
echo ""
#create entry for edge-hub 
$COMMAND_PATH/commands.sh edge-control-app edge-hub unix:uid:1000 &>> deploy-edge.log
check_command_status edge-app $?
echo "..Waiting for edge components to be initialized .." | tee -a deploy-edge.log
echo ""
sleep 5

echo "*Start edge SPIRE server for edge downstream component identity management (event-bus, user apps) " | tee -a deploy-edge.log
echo ""
#start downstream edge server- uses edge-hub certificates 
$COMMAND_PATH/commands.sh start-edge-server &>> deploy-edge.log
check_command_status start-edge-server $?
echo "..Waiting for edge components to be initialized .." | tee -a deploy-edge.log
echo ""
sleep 5

echo "*Start SPIRE agent edge downstream node which registers and attests edge downstream node to edge SPIRE server " | tee -a deploy-edge.log
echo ""
# register edge app node to edge server
$COMMAND_PATH/commands.sh edge-app-node &>> deploy-edge.log
check_command_status edge-app-node $?
echo "..Waiting for edge downstream server to be initialized .." | tee -a deploy-edge.log
echo ""
sleep 5

echo "*Start edge core components. Components include sidecars (Spiffe-helper - for certificate acquisition rotation, ghostunnel - for mTLS tunnel) " | tee -a deploy-edge.log
echo "*Workload attestation is performed against the unique selectors \"unix:uid:1000\"" | tee -a deploy-edge.log
echo "*All processes running under uid 1000 are attested and provided leaf certificates for communication with the help of spiffe-helper" | tee -a deploy-edge.log
echo ""
#create entry for event-bus
$COMMAND_PATH/commands.sh edge-app event-bus unix:uid:1000 &>> deploy-edge.log
check_command_status event-bus $?

#create entry for ligthmapper; this will require identifying unqiue process, run as different user id ? 
$COMMAND_PATH/commands.sh edge-app lightmapper unix:uid:1000 &>> deploy-edge.log
check_command_status lightmapper $?

#start edge-hub 
$COMMAND_PATH/commands.sh start-edge-hub &>> deploy-edge.log
check_command_status start-edge-hub $?

#start edge event-bus - helper and tunnel only
$COMMAND_PATH/commands.sh start-edge-event-bus &>> deploy-edge.log
check_command_status start-edge-event-bus $?

#start edge light mapper app - helper and tunnel only
$COMMAND_PATH/commands.sh start-edge-app &>> deploy-edge.log
check_command_status start-edge-app $?
echo "..Waiting for edge components to be initialized .." | tee -a deploy-edge.log
echo ""
sleep 5

echo "Edge security deployment is complete.Thank you." | tee -a deploy-edge.log
echo ""
spire=`ps -ef | grep spire` 
echo "$spire" | tee -a deploy-edge.log
cloud=`ps -ef | grep edge_core`
echo "$cloud" | tee -a deploy-edge.log
