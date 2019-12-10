//--------------------------------------------------------------------------------------------------
//
// Copyright (c) 2015-2019 Denis Dyakov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
// BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//--------------------------------------------------------------------------------------------------

package dht

// #include "gpio.h"
import "C"

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	"github.com/d2r2/go-shell"
	"github.com/davecgh/go-spew/spew"
)

// SensorType signify what sensor in use.
type SensorType int

// String implement Stringer interface.
func (v SensorType) String() string {
	if v == DHT11 {
		return "DHT11"
	} else if v == DHT12 {
		return "DHT12"
	} else if v == DHT22 || v == AM2302 {
		return "DHT22|AM2302"
	} else {
		return "!!! unknown !!!"
	}
}

// GetHandshakeDuration specify signal duration necessary
// to initiate sensor response.
func (v SensorType) GetHandshakeDuration() time.Duration {
	if v == DHT11 {
		return 18 * time.Millisecond
	} else if v == DHT12 {
		return 200 * time.Millisecond
	} else if v == DHT22 || v == AM2302 {
		return 18 * time.Millisecond
	} else {
		return 0
	}
}

// GetRetryTimeout return recommended timeout necessary
// to wait before new round of data exchange.
func (v SensorType) GetRetryTimeout() time.Duration {
	return 1500 * time.Millisecond
}

const (
	// DHT11 is most popular sensor
	DHT11 SensorType = iota + 1
	// DHT12 is more precise than DHT11 (has scale parts)
	DHT12
	// DHT22 is more expensive and precise than DHT11
	DHT22
	// AM2302 aka DHT22
	AM2302 = DHT22
)

// Pulse keep pulse signal state with how long it lasted.
type Pulse struct {
	Value    byte
	Duration time.Duration
}

// Activate sensor and get back bunch of pulses for further decoding.
// C function call wrapper.
func dialDHTxxAndGetResponse(pin int, handshakeDur time.Duration,
	boostPerfFlag bool) ([]Pulse, error) {

	var arr *C.int32_t
	var arrLen C.int32_t
	var list []int32
	var hsDurUsec C.int32_t = C.int32_t(handshakeDur / time.Microsecond)
	var boost C.int32_t = 0
	var err2 *C.Error
	if boostPerfFlag {
		boost = 1
	}
	// return array: [pulse, duration, pulse, duration, ...]
	r := C.dial_DHTxx_and_read(C.int32_t(pin), hsDurUsec,
		boost, &arr, &arrLen, &err2)
	if r == -1 {
		var err error
		if err2 != nil {
			msg := C.GoString(err2.message)
			err = errors.New(spew.Sprintf("Error during call C.dial_DHTxx_and_read(): %v", msg))
			C.free_error(err2)
		} else {
			err = errors.New(spew.Sprintf("Error during call C.dial_DHTxx_and_read()"))
		}
		return nil, err
	}
	defer C.free(unsafe.Pointer(arr))
	// convert original C array arr to Go slice list
	h := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	h.Data = uintptr(unsafe.Pointer(arr))
	h.Len = int(arrLen)
	h.Cap = int(arrLen)
	pulses := make([]Pulse, len(list)/2)
	// convert original int array ([pulse, duration, pulse, duration, ...])
	// to Pulse struct array
	for i := 0; i < len(list)/2; i++ {
		var value byte = 0
		if list[i*2] != 0 {
			value = 1
		}
		pulses[i] = Pulse{Value: value,
			Duration: time.Duration(list[i*2+1]) * time.Microsecond}
	}
	return pulses, nil
}

// decodeByte decode data byte from specific pulse array position.
func decodeByte(tLow, tHigh0, tHigh1 time.Duration, start int, pulses []Pulse) (byte, error) {
	if len(pulses)-start < 16 {
		return 0, fmt.Errorf("Can't decode byte, since range between "+
			"index and array length is less than 16: %d, %d", start, len(pulses))
	}
	HIGH_DUR_MAX := tLow + tHigh1
	HIGH_LOW_DUR_AVG := ((tLow+tHigh1)/2 + (tLow+tHigh0)/2) / 2
	var b int = 0
	for i := 0; i < 8; i++ {
		pulseL := pulses[start+i*2]
		pulseH := pulses[start+i*2+1]
		if pulseL.Value != 0 {
			return 0, fmt.Errorf("Low edge value expected at index %d", start+i*2)
		}
		if pulseH.Value == 0 {
			return 0, fmt.Errorf("High edge value expected at index %d", start+i*2+1)
		}
		// const HIGH_DUR_MAX = (70 + (70 + 54)) / 2 * time.Microsecond
		// Calc average value between 24us (bit 0) and 70us (bit 1).
		// Everything that less than this param is bit 0, bigger - bit 1.
		// const HIGH_LOW_DUR_AVG = (24 + (70-24)/2) * time.Microsecond
		if pulseH.Duration > HIGH_DUR_MAX {
			return 0, fmt.Errorf("High edge value duration %v exceed "+
				"maximum expected %v", pulseH.Duration, HIGH_DUR_MAX)
		}
		if pulseH.Duration > HIGH_LOW_DUR_AVG {
			//fmt.Printf("bit %d is high\n", 7-i)
			b = b | (1 << uint(7-i))
		}
	}
	return byte(b), nil
}

// Decode bunch of pulse read from DHTxx sensors.
// Use pdf specifications from /docs folder to read 5 bytes and
// convert them to temperature and humidity with results validation.
func decodeDHTxxPulses(sensorType SensorType, pulses []Pulse) (temperature float32,
	humidity float32, err error) {
	if len(pulses) >= 82 && len(pulses) <= 85 {
		pulses = pulses[len(pulses)-82:]
	} else {
		printPulseArrayForDebug(pulses)
		return -1, -1, fmt.Errorf("Can't decode pulse array received from "+
			"DHTxx sensor, since incorrect length: %d", len(pulses))
	}

	pulses = pulses[:80]
	// any bit low signal part
	tLow := 50 * time.Microsecond
	// 0 bit high signal part
	tHigh0 := 27 * time.Microsecond
	// 1 bit high signal part
	tHigh1 := 70 * time.Microsecond

	// decode 1st byte
	b0, err := decodeByte(tLow, tHigh0, tHigh1, 0, pulses)
	if err != nil {
		return -1, -1, err
	}
	// decode 2nd byte
	b1, err := decodeByte(tLow, tHigh0, tHigh1, 16, pulses)
	if err != nil {
		return -1, -1, err
	}
	// decode 3rd byte
	b2, err := decodeByte(tLow, tHigh0, tHigh1, 32, pulses)
	if err != nil {
		return -1, -1, err
	}
	// decode 4th byte
	b3, err := decodeByte(tLow, tHigh0, tHigh1, 48, pulses)
	if err != nil {
		return -1, -1, err
	}
	// decode 5th byte: control sum to verify all data received from sensor
	sum, err := decodeByte(tLow, tHigh0, tHigh1, 64, pulses)
	if err != nil {
		return -1, -1, err
	}
	// produce data consistency check
	calcSum := byte(b0 + b1 + b2 + b3)
	if sum != calcSum {
		err := errors.New(spew.Sprintf(
			"CRCs doesn't match: checksum from sensor(%v) != "+
				"calculated checksum(%v=%v+%v+%v+%v)",
			sum, calcSum, b0, b1, b2, b3))
		return -1, -1, err
	} else {
		lg.Debugf("CRCs verified: checksum from sensor(%v) = calculated checksum(%v=%v+%v+%v+%v)",
			sum, calcSum, b0, b1, b2, b3)
	}
	// debug output for 5 bytes
	lg.Debugf("Decoded from DHTxx sensor: [%d, %d, %d, %d, %d]", b0, b1, b2, b3, sum)
	// extract temperature and humidity depending on sensor type
	temperature, humidity = 0.0, 0.0
	if sensorType == DHT11 {
		humidity = float32(b0)
		temperature = float32(b2)
	} else if sensorType == DHT12 {
		humidity = float32(b0) + float32(b1)/10.0
		temperature = float32(b2) + float32(b3)/10.0
		if b3&0x80 != 0 {
			temperature *= -1.0
		}
	} else if sensorType == DHT22 {
		humidity = (float32(b0)*256 + float32(b1)) / 10.0
		temperature = (float32(b2&0x7F)*256 + float32(b3)) / 10.0
		if b2&0x80 != 0 {
			temperature *= -1.0
		}
	}
	// additional check for data correctness
	if humidity > 100.0 {
		return -1, -1, fmt.Errorf("Humidity value exceed 100%%: %v", humidity)
	} else if humidity == 0 {
		return -1, -1, fmt.Errorf("Humidity value cannot be zero")
	}
	// success
	return temperature, humidity, nil
}

// Print bunch of pulses for debug purpose.
func printPulseArrayForDebug(pulses []Pulse) {
	// var buf bytes.Buffer
	// for i, pulse := range pulses {
	// 	buf.WriteString(fmt.Sprintf("pulse %3d: %v, %v\n", i,
	// 		pulse.Value, pulse.Duration))
	// }
	// lg.Debugf("Pulse count %d:\n%v", len(pulses), buf.String())
	lg.Debugf("Pulses received from DHTxx sensor: %v", pulses)
}

// ReadDHTxx send activation request to DHTxx sensor via specific pin.
// Then decode pulses sent back with asynchronous
// protocol specific for DHTxx sensors.
//
// Input parameters:
// 1) sensor type: DHT11, DHT22 (aka AM2302);
// 2) pin number from GPIO connector to interact with sensor;
// 3) boost GPIO performance flag should be used for old devices
// such as Raspberry PI 1 (this will require root privileges).
//
// Return:
// 1) temperature in Celsius;
// 2) relative humidity in percent;
// 3) error if present.
func ReadDHTxx(sensorType SensorType, pin int,
	boostPerfFlag bool) (temperature float32, humidity float32, err error) {
	// activate sensor and read data to pulses array
	handshakeDur := sensorType.GetHandshakeDuration()
	pulses, err := dialDHTxxAndGetResponse(pin, handshakeDur, boostPerfFlag)
	if err != nil {
		return -1, -1, err
	}
	// output debug information
	printPulseArrayForDebug(pulses)
	// decode pulses
	temp, hum, err := decodeDHTxxPulses(sensorType, pulses)
	if err != nil {
		return -1, -1, err
	}
	return temp, hum, nil
}

// ReadDHTxxWithRetry send activation request to DHTxx sensor via specific pin.
// Then decode pulses sent back with asynchronous
// protocol specific for DHTxx sensors. Retry n times in case of failure.
//
// Input parameters:
// 1) sensor type: DHT11, DHT22 (aka AM2302);
// 2) pin number from gadget GPIO to interact with sensor;
// 3) boost GPIO performance flag should be used for old devices
// such as Raspberry PI 1 (this will require root privileges);
// 4) how many times to retry until success either counter is zeroed.
//
// Return:
// 1) temperature in Celsius;
// 2) relative humidity in percent;
// 3) number of extra retries data from sensor;
// 4) error if present.
func ReadDHTxxWithRetry(sensorType SensorType, pin int, boostPerfFlag bool,
	retry int) (temperature float32, humidity float32, retried int, err error) {
	// create default context
	ctx := context.Background()
	// reroute call
	return ReadDHTxxWithContextAndRetry(ctx, sensorType, pin,
		boostPerfFlag, retry)
}

// ReadDHTxxWithContextAndRetry send activation request to DHTxx sensor via specific pin.
// Then decode pulses sent back with asynchronous
// protocol specific for DHTxx sensors. Retry n times in case of failure.
//
// Input parameters:
// 1) parent context; could be used to manage life-cycle
//  of sensor request session from code outside;
// 2) sensor type: DHT11, DHT22 (aka AM2302);
// 3) pin number from gadget GPIO to interact with sensor;
// 4) boost GPIO performance flag should be used for old devices
//  such as Raspberry PI 1 (this will require root privileges);
// 5) how many times to retry until success either counter is zeroed.
//
// Return:
// 1) temperature in Celsius;
// 2) relative humidity in percent;
// 3) number of extra retries data from sensor;
// 4) error if present.
func ReadDHTxxWithContextAndRetry(parent context.Context, sensorType SensorType, pin int,
	boostPerfFlag bool, retry int) (temperature float32, humidity float32, retried int, err error) {
	// create context with cancellation possibility
	ctx, cancel := context.WithCancel(parent)
	// use done channel as a trigger to exit from signal waiting goroutine
	done := make(chan struct{})
	defer close(done)
	// build actual signals list to control
	signals := []os.Signal{os.Kill, os.Interrupt}
	if shell.IsLinuxMacOSFreeBSD() {
		signals = append(signals, syscall.SIGTERM)
	}
	// run goroutine waiting for OS termination events, including keyboard Ctrl+C
	shell.CloseContextOnSignals(cancel, done, signals...)
	retried = 0
	for {
		temp, hum, err := ReadDHTxx(sensorType, pin, boostPerfFlag)
		if err != nil {
			if retry > 0 {
				lg.Warning(err)
				retry--
				retried++
				select {
				// check for termination request
				case <-ctx.Done():
					// Interrupt loop, if pending termination.
					return -1, -1, retried, ctx.Err()
				// sleep before new attempt according to specification
				case <-time.After(sensorType.GetRetryTimeout()):
					continue
				}
			}
			return -1, -1, retried, err
		}
		return temp, hum, retried, nil
	}
}
