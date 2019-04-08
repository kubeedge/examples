#!/bin/bash

. ../edge.env

COMMAND_PATH=$SPIRE_PATH

check_command_status()
{
  if [ $2 -ne 0 ]; then
    echo $return
#TBD define exit codes for command failures
    exit 1
  else
    echo "command $1 is successful" >> demo/demo-scenario-2.log
  fi
  sleep 5
}

cleanapp()
{
	pkill spiffe-helper; pkill ghostunnel; sudo pkill edge_core;
# agent crashes with old certificates ??
	rm $SPIRE_PATH/user-app/certs/*
	rm $SPIRE_PATH/event-bus/certs/*
}

keypress()
{
  echo "Press any key to continue"
  read
}

cd $SPIRE_PATH

echo "**********************************************"
echo "Scenario-2 :[CLOUD OFFLINE SCENARIO]"
echo "Showcase after node and workload"
echo "attestation, cloud-side servers are not required"
echo "to issue or rotate certificates."
echo "Please note, there is a dependency to rotate CA"
echo "with the cloud server."
echo "**********************************************"

echo "**********************************************"
echo "[STEP-1] DEPLOY EDGE SECURITY AND COMPONENTS."
echo "**********************************************"
$COMMAND_PATH/deploy-edge.sh
sleep 5

echo "**********************************************"
echo "[STEP-2] SHUTDOWN CLOUD IDENTITY SERVER AND APPS"
echo "         EDGE APPS ARE RESTARTED TO DOWNLOAD NEW"
echo "         CERTIFICATES WITH CLOUD SERVER."
echo "**********************************************"
echo ""
echo "Make cloud component offline and press any key to continue"
echo "Execute \'deploy-cloud.sh clean\' command in cloud vm to shutdown all components in cloud"
echo ""
keypress
cleanapp

echo "Certifcates in $SPIRE_PATH/user-app/certs/"
ls $SPIRE_PATH/user-app/certs/
echo ""
echo "Certifcates in $SPIRE_PATH/event-bus/certs/"
ls $SPIRE_PATH/event-bus/certs/
echo ""
keypress


echo "******************************************************"
echo "[STEP-3] NEW CERTIFICATES ARE OBTAINED "
echo "         WITHOUT CLOUD IDENTITY SERVER - OFFLINE MODE"
echo "******************************************************"

echo "restart applications"

echo "*Start edge core components. Components include sidecars (Spiffe-helper - for certificate acquisition rotation, ghostunnel - for mTLS tunnel) " | tee -a demo/demo-scenario-2.log
echo "*Workload attestation is performed against the unique selectors \"unix:uid:1000\"" | tee -a demo/demo-scenario-2.log
echo "*All processes running under uid 1000 are attested and provided leaf certificates for communication with the help of spiffe-helper" | tee -a demo/demo-scenario-2.log
echo ""

#start edge-hub 
$COMMAND_PATH/commands.sh start-edge-hub &>> demo/demo-scenario-2.log
check_command_status start-edge-hub $?

#start edge event-bus - helper and tunnel only
$COMMAND_PATH/commands.sh start-edge-event-bus &>> demo/demo-scenario-2.log
check_command_status start-edge-event-bus $?

#start edge light mapper app - helper and tunnel only
$COMMAND_PATH/commands.sh start-edge-app &>> demo/demo-scenario-2.log
check_command_status start-edge-app $?
echo "..Waiting for edge components to be initialized .." | tee -a demo/demo-scenario-2.log
echo ""
sleep 5

echo "Certifcates in $SPIRE_PATH/user-app/certs/"
ls $SPIRE_PATH/user-app/certs/
echo ""
echo "Certifcates in $SPIRE_PATH/event-bus/certs/"
ls $SPIRE_PATH/event-bus/certs/
echo ""

cd $SPIRE_PATH/demo
