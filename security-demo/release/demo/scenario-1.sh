#!/bin/bash

. /opt/spire/edge.env

COMMAND_PATH=/opt/spire

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

keypress()
{
  echo "Press any key to continue"
  read
}

echo "**********************************************"
echo "Scenario-1 : Showcase without node attestation"
echo "and workload attestation, no certificates/svid"
echo "are downloaded for the workloads"
echo "**********************************************"

cd /opt/spire

date=`date`
echo "**Started at $date" | tee -a demo/demo-scenario-1.log
echo ""

echo "*Generating configurations" | tee -a demo/demo-scenario-1.log
echo ""
$COMMAND_PATH/generate-edge-config.sh

echo "*********************************************"
echo "[STEP-1] REGISTER EDGE NODE TO UPSTREAM SPIRE"
echo "         SERVER. NODE IS ATTESTED BASED ON"
echo "         TOKEN METHOD."
echo "*********************************************"

echo "*Start SPIRE agent at edge node which registers and attests edge node to cloud " | tee -a demo/demo-scenario-1.log
echo ""
# register edge node to upstream server
$COMMAND_PATH/commands.sh edge-node &>> demo/demo-scenario-1.log
check_command_status edge-node $?
echo "..Waiting for edge node to be initialized .." | tee -a demo/demo-scenario-1.log
echo ""
sleep 5

echo "*********************************************"
echo "[STEP-2] VERIFY NODE ENTRIES, NOTE THERE IS "
echo "         NO ENTRY FOR EDGE HUB(EDGE CORE APP)."
echo "*********************************************"
echo ""

entry=$(sshpass -p $CLOUD_VM_PASS ssh $CLOUD_VM_USER@$CLOUD_VM_IP $SPIRE_PATH/spire-server entry show 3>&1)
err=$?
echo $entry

keypress


echo "*********************************************"
echo "[STEP-3] START EDGE CORE WITH REGISTRATION TO"
echo "         SPIRE SERVER."
echo "*********************************************"
echo ""

echo "***Start edge core without entry.."
#start edge-hub 
$COMMAND_PATH/commands.sh start-edge-hub &>> deploy-edge.log
check_command_status start-edge-hub $?

echo "Let's checkout the logs - helper.log"
echo ""

echo "Check certificates"
echo ""

echo "*********************************************"
echo "[STEP-4] CERTIFICATES ARE NOT DOWNLOADED."
echo "*********************************************"
echo ""

ls $COMMAND_PATH/certs/svid.pem 

echo "Check helper log"
cat ./log/helper.log
keypress

#kill the previously started sidecars and app
sudo pkill edge_core; pkill spiffe-helper; pkill ghostunnel

echo "*********************************************"
echo "[STEP-5] START EDGE CORE WITH REGISTRATION."
echo "         EDGE CORE IS REGISTERED AND ATTESTED."
echo "*********************************************"
echo ""

echo "***Create entry for edge core"

echo "*Start edge core components. Components include sidecars (Spiffe-helper - for certificate acquisition rotation, ghostunnel - for mTLS tunnel) " | tee -a demo/demo-scenario-1.log
echo "*Workload attestation is performed against the unique selectors \"unix:uid:1000\"" | tee -a demo/demo-scenario-1.log
echo "*All processes running under uid 1000 are attested and provided leaf certificates for communication with the help of spiffe-helper" | tee -a demo/demo-scenario-1.log
echo ""
#create entry for edge-hub 
$COMMAND_PATH/commands.sh edge-control-app edge-hub unix:uid:1000 &>> deploy-edge.log
check_command_status edge-app $?
echo "..Waiting to create entry.." | tee -a deploy-edge.log
echo ""
sleep 5


#start edge-hub 
$COMMAND_PATH/commands.sh start-edge-hub &>> demo/demo-scenario-1.log
check_command_status start-edge-hub $?

echo "..Waiting for edge components to be initialized .." | tee -a demo/demo-scenario-1.log
echo ""
sleep 5

echo "Let's checkout the logs - helper.log"
echo ""

echo "*********************************************"
echo "[STEP-6] CERTIFICATES ARE NOT DOWNLOADED."
echo "*********************************************"
echo ""

echo "Check certificates"
echo ""

ls $COMMAND_PATH/certs/svid.pem 
echo""

echo "Check helper log"
cat ./log/helper.log
keypress

cd /opt/spire/demo
