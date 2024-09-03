package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/goburrow/modbus"
	"k8s.io/klog/v2"
)

const TOTAL_SIZE int32 = 320044 // expected file size
const CHUNK_SIZE uint16 = 60000 // number of registers to read each time

const BLOCKSIZE uint16 = 100 // number of registers to read each time
const ISREAD uint16 = 1      // read flag
const NOREAD uint16 = 0      // not read flag

// Get notification from modbus server
func getNotification(client modbus.Client, chunk_id uint16) bool {
	results, err := client.ReadHoldingRegisters(chunk_id, 1)
	if err != nil {
		klog.Fatal(err)
	}

	if len(results) == 0 {
		klog.Fatal("No data received")
	}
	// Convert string to byte slice
	bytes := []byte(results)
	// Convert the byte slice to uint16 using binary.BigEndian
	intResults := binary.BigEndian.Uint16(bytes)
	if intResults == NOREAD {
		return true
	}
	return false
}

// Send notification to modbus server
func sendNotification(client modbus.Client, chunk_id uint16) {
	_, err := client.WriteSingleRegister(chunk_id, ISREAD) // Write to register
	if err != nil {
		klog.Fatal(err)
	}
}

// Get the modbus client
func getClient() modbus.Client {
	handler := modbus.NewTCPClientHandler("localhost:5020") // Use TCP communication
	handler.Timeout = 1 * 1e9                               // 1 second
	handler.SlaveId = 1
	err := handler.Connect()
	if err != nil {
		klog.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	return client
}

// Save the audio file
func saveWavFile(audioData *bytes.Buffer, name string) {
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "tempfile-*.wav")
	if err != nil {
		klog.Errorf("Error creating temp file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name()) // Make sure the temporary file is deleted at the end of the function

	// Write the contents of audioData to a temporary file
	_, err = tempFile.Write(audioData.Bytes())
	if err != nil {
		klog.Errorf("Error writing to temp file: %v", err)
		return
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		klog.Errorf("Error closing temp file: %v", err)
		return
	}

	cmd := exec.Command("mv", tempFile.Name(), name)
	err = cmd.Run()
	if err != nil {
		klog.Fatalf("Failed to move file: %v", err)
		return
	}

	klog.Infoln("Save wav file successfully")
}

// Get the minimum value of two int32
func min(a, b int32) uint16 {
	if a < b {
		return uint16(a)
	}
	return uint16(b)
}

// Receive audio data
func receiveAudioData(client modbus.Client, chunk_id uint16, audioData *bytes.Buffer) {
	// Calculate the number of registers read by receiveAudioData
	tmp := TOTAL_SIZE/2 - int32(chunk_id)*int32(CHUNK_SIZE)
	rec_size := min(int32(CHUNK_SIZE), tmp)

	// Read
	var startAddr uint16 = 0
	for ; ; startAddr += BLOCKSIZE {
		// If the number of registers read is less than or equal to 0, exit
		if rec_size <= startAddr {
			break
		}

		// Calculate the number of registers read each time
		var singleSize = BLOCKSIZE
		if startAddr+BLOCKSIZE > rec_size {
			singleSize = rec_size - startAddr
		}

		results, err := client.ReadInputRegisters(uint16(startAddr), uint16(singleSize))
		if err != nil {
			klog.Fatal(err)
		}
		if len(results) == 0 {
			break
		}

		// Write the result to audioData
		for i := 0; i < len(results); i += 2 {
			audioData.WriteByte(results[i])
			if i+1 < len(results) {
				audioData.WriteByte(results[i+1])
			}
		}
	}
}

func main() {
	const INTVAL = 1 //Interval time
	client := getClient()

	chunkSize := int32(CHUNK_SIZE) * 2
	numChunks := math.Ceil(float64(TOTAL_SIZE) / float64(chunkSize))
	NUM_CHUNK := uint16(numChunks)

	for {
		startTime := time.Now()
		var audioData bytes.Buffer
		var i uint16 = 0
		for ; i < NUM_CHUNK; i++ {
			for {
				if getNotification(client, uint16(i)) {
					receiveAudioData(client, uint16(i), &audioData)
					sendNotification(client, uint16(i))
					break
				}
			}
		}
		saveWavFile(&audioData, "/etc/data/received.wav")
		elapsedTime := time.Since(startTime) // Calculate elapsed time
		if elapsedTime < INTVAL {
			time.Sleep(INTVAL - elapsedTime) // Sleep to ensure total time per file is INTVAL seconds
		}
	}
}
