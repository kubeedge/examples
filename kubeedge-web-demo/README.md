# KubeEdge Web Demo

Firstly the users open a browser,
and enter the web app page by the web app link,
choose the music and click the button `Play` in the web page,
at last the expected track is pushed to the edge node
and the track is played on the speaker connected to the edge node.

## Prerequisites

### Hardware Prerequisites

* RaspBerry PI (RaspBerry PI 3 has been used for this demo).
  The RaspBerry PI is also the edge node to which the speaker will be connected.

* A speaker for playing the music.

### Software Prerequisites

* A running Kubernetes cluster.

* KubeEdge [1.0.0](https://github.com/kubeedge/kubeedge/releases/tag/v1.0.0)
  See [instructions](https://github.com/kubeedge/kubeedge/blob/master/docs/getting-started/usage.md#run-kubeedge) on how to setup KubeEdge.

  *Note*:

  when you setup edgecore on the RaspBerry PI,
  Please set the `mqtt mode` into `2` [in this line](https://github.com/kubeedge/kubeedge/blob/master/edge/conf/edge.yaml#L4),
  and replace `0.0.0.0` with your Kubernetes master ip address [in this line](https://github.com/kubeedge/kubeedge/blob/master/edge/conf/edge.yaml#L11).

* In order to control the speaker and play the expected track, we need to manage the speaker connected to the RaspBerry PI.
  KubeEdge allows us to manage devices using Kubernetes custom resource definitions.
  The design proposal is [here](https://github.com/kubeedge/kubeedge/blob/master/docs/proposals/device-crd.md).
  Apply the CRD schema yamls available [here](https://github.com/kubeedge/kubeedge/tree/master/build/crds/devices) using kubectl. 

## Steps to run the demo

### Clone demo code

```console
 git clone https://github.com/kubeedge/examples
```

### Create the device model and device instance for the speaker

With the Device CRD APIs now installed in the cluster,
we create the device model and instance for the speaker using the yaml files.

```console
 cd $GOPATH/src/github.com/kubeedge/examples/kubeedge-web-demo/kubeedge-web-app/deployments/
 kubectl create -f kubeedge-speaker-model.yaml
 kubectl create -f kubeedge-speaker-instance.yaml
```

### Run KubeEdge Web App

The KubeEdge Web App runs in a VM on cloud.
It can be deployed using a Kubernetes deployment yaml.

```console
 cd $GOPATH/src/github.com/kubeedge/examples/kubeedge-web-demo/kubeedge-web-app/deployments/
 kubectl create -f kubeedge-web-app.yaml
```

### Build PI Player App

Cross-complie the PI Player App which will run on the RaspBerry PI and play the expected track.

```console
 cd $GOPATH/src/github.com/kubeedge/examples/kubeedge-web-demo/pi-player-app/
 export GOARCH=arm
 export GOOS="linux"
 export GOARM=6
 export CGO_ENABLED=1
 export CC=arm-linux-gnueabi-gcc
 go build -o pi-player-app main.go
```

### Run PI Player App

Make sure the MQTT broker is running on the RaspBerry PI.
Copy the PI Player App binary to the RaspBerry PI and run it.
The App will subscribe to the `$hw/events/device/speaker-01/twin/update/document` topic 
and when it receives the expected track on the topic, it will play it on the speaker.
At last, you need to copy the music files into the folder `/home/pi/music/` on the RaspBerry PI.
The music file name is like <track>.mp3, for example: `1.mp3`

The PI Player App will issue the `omxplayer` to play music,
please make sure the `omxplayer` is installed on the RaspBerry PI.
If not, please see the following link to setup `omxplayer`.

https://raspberry-projects.com/pi/software_utilities/media-players/omxplayer

```console
 ./pi-player-app
```

### Play music by visiting Web App Page

* Visit web app page by the web app link.

* Choose the music you want to play, and then click the button `Play`.
  The track info is pushed to the RaspBerry PI and the music is played on the speaker.

* Click the button `stop` to stop the music, the music is stopped on the speaker.
