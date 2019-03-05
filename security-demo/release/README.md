#Identity Management Using SPIRE for Kubeedge

##SPIRE-based Identity Management Reference Architecture

Goal of this example is to achieve workload identity management for
cloud and edge processes . A reference architecture is shown below . A
list of things to-do is listed below in version information section .
For further integration , an extension add-on will be developed
integrated with kubeege for the same .

![](doc/images/demo_arch.png)

##Folders
* release : Directory has configurations and scripts to be used for deployment of identity management infrastructure for Kubeedge.

* app-agent-conf  : Configurations for spire agent interfacing with edge (event-bus) and user applications . Spire agent communicates to edge spire server.|

* conf : Configurations for spire server (cloud and edge) and spire agent (agent communicating to cloud spire server).

* certs : Spiffe helper communicates with spire agents and download certificates in this folder . Ghostunnel is loaded with these certificates . For edge to user interface application , spiffe helper configuration and ghostunnel certs are present in event-bus and user-app folders respectively.

* app-binaries : All binaries for edge , cloud and test applications.

* log : Logs for different processes run as part of the example.

##Scripts
* deploy-cloud.sh : Deploys cloud part of identity management infrastructure. Used to 
                    
                    1) Start cloud spire server.
                           
                    2) Registers cloud node and cloud agent.
                           
                    3) Registers and starts cloud hub (test cloud hub command simulator).

* deploy-edge.sh  : Deploys edge part of identity management infrastructure.
                           
                    1) Registers edge node and edge agent communication with cloud spire
                           server.
                           
                    2) Registers and starts edge spire server.
                           
                    3) Registers and starts edge agent for edge to user app communication interface.

* start-spiffe-helper.sh  : Based on the configuration (helper.conf) , starts spiffe-helper communication for certificate download and rotation. Ghostunnel is run using spiffe-helper.

* commands.sh  : Abstracts spire cli commands.
 
##Configuration
### Cloud node configuration

*Upstream CA – Cloud spire server configuration:*
&lt;SPIRE\_PATH&gt;/conf/server/server.conf

*Upstream CA Agent – Cloud node agent configuration :*
&lt;SPIRE\_PATH&gt;/conf/agent/agent.conf

Following fields might require modification based on the deployment environment

       bind_address = "192.168.56.101" <Cloud node vm IP>

       server_address = "192.168.56.101" <Cloud Spire Server IP>

       server_port = "8081" <Cloud Spire Server port>

       socket_path =/tmp/agent.sock"

### Edge node configuration

*Edge agent configuration to **connect to cloud spire server** :*
&lt;SPIRE\_PATH&gt;/conf/agent/agent.conf

Following fields might require modification based on the deployment environment

    bind_address = "192.168.56.102"

    bind_port = "8088"

    data_dir = "./.data"

    server_address = "192.168.56.101" // <Cloud spire server IP>

    server_port = "8081" // <Cloud spire server port>

    socket_path ="/tmp/upstream_agent.sock"

*Edge spire server configuration :*
&lt;SPIRE\_PATH&gt;/conf/server/server.conf

    UpstreamCA "spire" {

    plugin_data {

      server_address = "192.168.56.101" // <Cloud spire server IP>

      server_port = 8081 // <Cloud spire server port>

      workload_api_socket = "/tmp/upstream\_agent.sock"

    }

*Edge application agent configuration to connect to edge spire server :
&lt;SPIRE\_PATH&gt;/conf/app-agent-conf/agent/agent.conf*

    bind_address = "192.168.56.102"

    bind_port = "9088"

    data_dir = "./.app-data"

    server_address = "192.168.56.102" // <Edge spire server IP>

    server_port = "8081" // <Edge spire server port>

    socket_path ="/tmp/app-agent.sock"


### Spiffe helper configuration (IMPORTANT)

Spiffe helper is used to execute ghostunnel for creating communciation
channel over TLS for cloud , edge and user workloads . Please refer to
the specific configurations in the folders section.

All the helper configurations need to updated with IP addresses and
ports present in

1\. &lt;SPIRE\_PATH&gt;/helper.conf

2\. &lt;SPIRE\_PATH&gt;/event-bus/event-bus-helper.conf

3\. &lt;SPIRE\_PATH&gt;/user-app/user-app-helper.conf

In the list of To-Dos , there is an item to generate these
configurations automatically based on environment configuration to avoid
errors.

Following is a sample configuration for kubeedge event-bus interface:

agentAddress = "/tmp/app-agent.sock"

cmdArgs = "**server** --listen 192.168.56.102:18884 --target
127.0.0.1:2884 --cacert /opt/spire/event-bus/certs/svid\_bundle.pem
--keystore /opt/spire/event-bus/certs/new\_bundle.p12 **--allow-uri-san
spiffe://example.org/downstream-app-event-bus**"

cmd = "/opt/spire/ghostunnel"

certDir = "/opt/spire/event-bus/certs"

renewSignal = "SIGUSR1"

svidFileName = "svid.pem"

svidKeyFileName = "svid\_key.pem"

svidBundleFileName = "svid\_bundle.pem"

### Script Configurations

Environment variable configurations : &lt;SPIRE\_PATH&gt;/edge.env

export CLOUD\_VM\_USER=vm1

export CLOUD\_VM\_PASS=vm1

export CLOUD\_VM\_IP=192.168.56.101

export SPIRE\_PATH=/opt/spire

Usage
-----

1\. Copy the package to cloud node and edge node .

2\. Update the IP , port and spire path in the configurations.

3\. In cloud node , execute deploy-cloud.sh.

4\. In edge node , execute deploy-edge.sh.

5\. Register an example device with cloud using following command. Please
note, in the current version, cloud test application opens 30000 port
for metadata creation (create pod or device) and 20000 port for
communication with kubeedge edgehub.

curl -XGET http://127.0.0.1:30000/device -H 'content-type:
application/json' -d@&lt;SPIRE\_PATH&gt;/app-binaries/test-device.yaml

6\. Run the light\_mapper application from app\_binaries. Light mapper
application is a binary built from
[*https://github.com/kubeedge/examples/tree/master/led-raspberrypi*](https://github.com/kubeedge/examples/tree/master/led-raspberrypi).
Usage of the application can be referred in the same page.

For spire server cli usage , please refer to
[*https://github.com/spiffe/spire*](https://github.com/spiffe/spire).

Present support
---------------

1\. Initial infrastructure deployment.

2\. 1-1 cloud and edge node. Multiple edge nodes are not supported in
scripts.

3\. Communication from user application and edge application using
certificates issued by cloud spire server and edge spire server.

ToDo
----

1\. Optimization for redundancies in configuration and scripts.

2\. Automate build/release and remove binaries from the package.

3\. Test and support upstream\_bundle=false for edge spire server.

4\. Certificate rotation issue needs to be automated or changes may be
required in spiffe helper for the dependency on keystore for ghostunnel.

5\. Auto-generate configurations (partially) based on environment
information.
