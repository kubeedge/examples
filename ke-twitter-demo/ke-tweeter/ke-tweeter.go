package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"k8s.io/client-go/rest"
	"github.com/ke-twitter-demo/ke-tweeter/constants"
	
	"github.com/ke-twitter-demo/ke-tweeter/utils"
	"github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller/apis/devices/v1alpha1"
)

// DeviceStatus is used to patch device status
type DeviceStatus struct {
	Status v1alpha1.DeviceStatus `json:"status"`
}

// The device id of the speaker
var deviceID = "speaker-01"

// The default namespace in which the speaker device instance resides
var namespace = "default"

// A regular expression which evaluates to true if the tweeted text
// contains `kubeedge play`.
var r = regexp.MustCompile(`\w*kubeedge play\w*`)

// The CRD client used to patch the device instance.
var crdClient *rest.RESTClient

// Twitter credentials required to build the Twitter client.
var accessToken, accessTokenSecret, consumerKey, consumerSecret string

func main() {

	// Create a client to talk to the K8S API server to patch the device CRDs
	kubeConfig, err := utils.KubeConfig()
	if err != nil {
		log.Fatalf("Failed to create KubeConfig , error : %v", err)
	}
	crdClient, err = utils.NewCRDClient(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create device crd client , error : %v", err)
	}

	// Read twitter credentials from the secret mounted in the container.
	err = readTwitterCredentialsFromSecret()

	twitterCredentials := utils.Credentials{
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
	}
	// Create a Twitter client passing the credentials built above
	client, err := utils.GetTwitterClient(&twitterCredentials)
	if err != nil {
		log.Fatalf("Error creating Twitter Client : %v", err)
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweetText := tweet.Text
		log.Printf("Got tweet %s\n", tweetText)
		tweetText = strings.TrimSpace(tweetText)
		tweetText = strings.ToLower(tweetText)

		if r.MatchString(tweetText) {
			items := r.Split(tweetText, -1)
			songTrack := strings.TrimSpace(items[1])
			UpdateDeviceTwinWithDesiredTrack(songTrack)
		} else {
			log.Println("Couldn't parse tweet, please tweet in this format:  [ kubeedge play xyz ]")
		}
	}
	log.Println("Start watching KubeEdge tweets")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"kubeedge", "Kubeedge", "KubeEdge", "KUBEEDGE"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Println("Stop watching KubeEdge tweets")
	stream.Stop()
}

// readTwitterCredentialsFromSecret uses the utility functions
// to read the credentials to authorize the application with Twitter.
// The credentials are assumed to be mounted via a Kubernetes secret
// inside the container.
func readTwitterCredentialsFromSecret() error {
	var err error
	accessToken, err = utils.ReadSecretKey(constants.TwitterAccessToken)
	if err != nil {
		log.Fatalf("Error reading %s : %v", constants.TwitterConsumerKey, err)
	}
	accessTokenSecret, err = utils.ReadSecretKey(constants.TwitterAccessTokenSecret)
	if err != nil {
		log.Fatalf("Error reading %s : %v", constants.TwitterConsumerKey, err)
	}
	consumerKey, err = utils.ReadSecretKey(constants.TwitterConsumerKey)
	if err != nil {
		log.Fatalf("Error reading %s : %v", constants.TwitterConsumerKey, err)
	}
	consumerSecret, err = utils.ReadSecretKey(constants.TwitterConsumerSecret)
	if err != nil {
		log.Fatalf("Error reading %s : %v", constants.TwitterConsumerSecret, err)
	}
	return err
}

// UpdateDeviceTwinWithDesiredTrack patches the desired state of
// the device twin with the track to play.
func UpdateDeviceTwinWithDesiredTrack(track string) bool {
	status := buildStatusWithDesiredTrack(track)
	deviceStatus := &DeviceStatus{Status: status}
	body, err := json.Marshal(deviceStatus)
	if err != nil {
		log.Fatalf("Failed to marshal device status %v", deviceStatus)
		return false
	}
	result := crdClient.Patch(constants.MergePatchType).Namespace(namespace).Resource(constants.ResourceTypeDevices).Name(deviceID).Body(body).Do()
	if result.Error() != nil {
		log.Fatalf("Failed to patch device status %v of device %v in namespace %v \n error:%+v", deviceStatus, deviceID, namespace, result.Error())
		return false
	} else {
		log.Printf("Track [ %s ] will be played on speaker %s", track, deviceID)
	}
	return true
}

func buildStatusWithDesiredTrack(song string) v1alpha1.DeviceStatus {
	metadata := map[string]string{"timestamp": strconv.FormatInt(time.Now().Unix()/1e6, 10),
		"type": "string",
	}
	twins := []v1alpha1.Twin{{PropertyName: "track", Desired: v1alpha1.TwinProperty{Value: song, Metadata: metadata}}}
	devicestatus := v1alpha1.DeviceStatus{Twins: twins}
	return devicestatus
}
