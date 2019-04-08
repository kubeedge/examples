#!/bin/bash 

COMMAND_PATH=/opt/spire

clean()
{
  pkill spire-server; pkill spire-agent; pkill spiffe-helper; pkill ghostunnel; pkill cloud-app; pkill edgecontroller
  rm .data/*
  rm .agentdata/*
  rm certs/*
  cp /opt/spire/conf/server/server.conf.bk /opt/spire/conf/server/server.conf
}

check_command_status()
{
  if [ $2 -ne 0 ]; then
    echo $return
#TBD define exit codes for command failures
    exit 1
  else 
    echo "command $1 is successful" >> deploy-cloud.log
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

date=`date`
echo "**Started at $date" | tee -a deploy-cloud.log
echo ""

echo "*Generating configurations" | tee -a deploy-cloud.log
echo ""
$COMMAND_PATH/generate-cloud-config.sh

echo "*Start SPIRE server at cloud node" | tee -a deploy-cloud.log
echo ""
$COMMAND_PATH/commands.sh start-cloud-server &>> deploy-cloud.log
check_command_status start-server $?

echo "..Waiting for server to be initialized ..." | tee -a deploy-cloud.log
echo ""
sleep 5

echo "*Start SPIRE agent for cloud which registers and attests the cloud node" | tee -a deploy-cloud.log
echo ""
$COMMAND_PATH/commands.sh cloud-node &>> deploy-cloud.log
check_command_status cloud-node $?

echo "..Waiting for node to be initialized ..." | tee -a deploy-cloud.log
echo ""
sleep 5

echo "*Start cloud components. Components include sidecars (Spiffe-helper - for certificate acquisition rotation, ghostunnel - for mTLS tunnel)" | tee -a deploy-cloud.log
echo "*Workload attestation is performed against the unique selectors \"unix:uid:1000\"" | tee -a deploy-cloud.log
echo "*All processes running under uid 1000 are attested and provided leaf certificates for communication with the help of spiffe-helper" | tee -a deploy-cloud.log
echo ""

$COMMAND_PATH/commands.sh cloud-app cloud-hub unix:uid:1000 &>> deploy-cloud.log
check_command_status cloud-app $?

$COMMAND_PATH/commands.sh start-cloud-hub &>> deploy-cloud.log
check_command_status start-cloud-hub $?

echo "..Waiting for hub initialization ..." | tee -a deploy-cloud.log
echo ""
sleep 5

echo "Cloud security deployment is complete.Thank you." | tee -a deploy-cloud.log

echo ""
spire=`ps -ef | grep spire` 
echo "$spire" | tee -a deploy-cloud.log
cloud=`ps -ef | grep edgecontroller`
echo "$cloud" | tee -a deploy-cloud.log
