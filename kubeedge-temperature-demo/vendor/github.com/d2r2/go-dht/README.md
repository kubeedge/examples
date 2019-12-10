DHTxx temperature and humidity sensors
======================================

[![Build Status](https://travis-ci.org/d2r2/go-dht.svg?branch=master)](https://travis-ci.org/d2r2/go-dht)
[![Go Report Card](https://goreportcard.com/badge/github.com/d2r2/go-dht)](https://goreportcard.com/report/github.com/d2r2/go-dht)
[![GoDoc](https://godoc.org/github.com/d2r2/go-dht?status.svg)](https://godoc.org/github.com/d2r2/go-dht)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
<!--
[![Coverage Status](https://coveralls.io/repos/d2r2/go-dht/badge.svg?branch=master)](https://coveralls.io/r/d2r2/go-dht?branch=master)
-->

Preamble
--------

All sensors supported by this library communicate with Raspberry PI via "single-wire digital interface". This approach is require very accurate synchronous timing. Raspberry PI hardware and its clones equipped with Linux - most massive and useful operating system (Linux is my favorite OS) - but not a real time system strictly speaking. It cannot guaranty that code (particularly user, not kernel layer code) might provide exact microsecond switching necessary to meet "single-wire digital interface" specification. Especially when you run this code from language whose runtime is based on background garbage collector processes, which might additionally raise "stop the world" phenomenon, when running code might unexpectedly freeze for a moment. Starting from Go 1.5 GC STW significantly improved, but still can affect millisecond/microsecond world.

So here is my recommendation - if you want to control sensors and peripheral devices from Golang environment running on Linux OS, go and buy sensors with **I2C interface**. I personally support such sensors in my hobby projects via [I2C library](https://github.com/d2r2/go-i2c) (you can go and find list of supported devices in the middle).

But if you want to teach how signal processing algorithm attempt to work with "sinlge-wire digital interface" from not real time user code layer, go and find out!

About
-----

DHT11 ([pdf reference](https://raw.github.com/d2r2/go-dht/master/docs/DHT11.pdf)), AM2302/DHT22 ([pdf reference](https://raw.github.com/d2r2/go-dht/master/docs/AM2302.pdf)) and DHT12 ([pdf reference](https://raw.github.com/d2r2/go-dht/master/docs/DHT12.pdf)) sensors, which quite popular among Arduino, Raspberry PI developers (here you will find comparison [DHT11 vs DHT22](https://raw.github.com/d2r2/go-dht/master/docs/dht.pdf)):
![dht11 and dht22](https://raw.github.com/d2r2/go-dht/master/docs/dht11_dht22.jpg)

They are cheap enough and affordable. So, here is a code written in [Go programming language](https://golang.org/) for Raspberry PI and clones, which gives you at the output temperature and humidity values (making all necessary signal processing via their own 1-wire bus protocol behind the scene).


Technology overview
-------------------

There are 2 methods how we can drive such devices which require special pins switch from low to high level and back (employing specific "1-wire protocol" described in pdf documentation):
1) First approach implies to work on the most lower layer to handle pins via GPIO chip registers using Linux "memory mapped" device (/dev/mem). This approach is most reliable (until you move to other RPI clone) and fastest with regard to the transmission speed. Disadvantage of this method is explained by the fact that each RPI-clone have their own GPIO registers set to drive device GPIO pins.
2) Second option implies to access GPIO pins via special layer based on Linux "device tree" approach (/sys/class/gpio/... virtual file system), which translate such operations to direct register writes and reads described in 1st approach. In some sense it is more compatible when you move from original Raspberry PI to RPI-clones, but may have some issues in stability of specific implementations. As it was found some clones don't implement this layer at all from the box (Beaglebone for instance). 

So, here I'm using second approach.

Compatibility
-------------

Tested on Raspberry PI 1/2 (model B), Banana PI (model M1), Orange PI One.

Golang usage
------------

```go
func main() {
	// Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
	// You may enable "boost GPIO performance" parameter, if your device is old
	// as Raspberry PI 1 (this will require root privileges). You can switch off
	// "boost GPIO performance" parameter for old devices, but it may increase
	// retry attempts. Play with this parameter.
	// Note: "boost GPIO performance" parameter is not work anymore from some
	// specific Go release. Never put true value here.
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(dht.DHT11, 4, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	// Print temperature and humidity
	fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
		temperature, humidity, retried)
}
```

Installation
------------

```bash
$ go get -u github.com/d2r2/go-dht
```

Quick start
-----------

There are two functions you could use: ```ReadDHTxx(...)``` and ```ReadDHTxxWithRetry(...)```.
They both do exactly same thing - activate sensor then read and decode temperature and humidity values.
The only thing which distinguish one from another - "retry count" parameter as additional argument in ```ReadDHTxxWithRetry(...)```.
So, it's highly recommended to utilize ```ReadDHTxxWithRetry(...)``` with "retry count" not less than 7, since sensor asynchronous protocol is not very stable causing errors time to time. Each additional retry attempt takes 1.5-2 seconds (according to specification before repeated attempt you should wait 1-2 seconds).

This functionality works not only with Raspberry PI, but with counterparts as well (tested with Raspberry PI and Banana PI).

> Note: This package does not have dependency on any sensor-specific 3-rd party C-code or library.

Tutorial
--------

The library consists of 2 parts: low level C-code to send queries and read raw data from sensor and front-end Golang functions with raw data decoding.

Originally attempt was made to write whole library in Golang, but during debugging it was found that Garbage Collector (GC) "stop the world" characteristic in early version of Golang noticeably freeze library in the middle of sensor reading process, which lead to unpredictable mistakes when some signals from sensor just have been lost.  Starting from Go 1.5 version GC behavior was improved significantly, but original design left as is, since it been tested and works reliably in most cases.

To install library on your Raspberry PI device you should execute console command `go get -u github.com/d2r2/go-dht` to download and install/update package to you device `$GOPATH/src` path.

You may start from simple test with DHTxx sensor using `./examples/example1/example1.go` application which will interact with the sensor connected to some specific physical hardware pin (you may google pinout of any Raspberry PI version either its clones).

Also you can use cross compile technique, to build ARM application from x86/64bit system. For this your should install GCC tool-chain for ARM target platform. So, your x86/64bit Linux system should have specific gcc compiler installed: in case of Debian or Ubuntu `arm-linux-gnueabi-gcc` (in case of Arch Linux `arm-linux-gnueabihf-gcc`).
After all, for instance, for cross compiling test application "./examples/example1/example1.go" to ARM target platform in Ubuntu/Debian you should run `CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build ./examples/example1/example1.go`.

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-dht).

For detailed explanation read great article "[Golang with Raspberry Pi : Read RH and Temperature from DHT22 or AM2302](https://skylabin.wordpress.com/2015/09/18/golang-with-raspberry-pi-read-rh-and-temperature-from-dht22-or-am2302)" written by Joseph Mathew. Thanks Joseph!

Contribute authors
------------------

* Joseph Mathew (https://skylabin.wordpress.com/)
* Alex Zhang ([ztc1997](https://github.com/ztc1997))
* Andy Brown ([andybrown668](https://github.com/andybrown668))
* Gareth Dunstone ([gdunstone](https://github.com/gdunstone))

Contact
-------

Please use [Github issue tracker](https://github.com/d2r2/go-dht/issues) for filing bugs or feature requests.

License
-------

Go-dht is licensed under MIT License.
