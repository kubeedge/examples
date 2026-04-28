package counter

import (
	"log"
	"time"
)

const (
	ON = iota
	OFF
)

type Counter struct {
	status chan int
	handle func (int)
}


func (counter *Counter) runDevice(interrupt chan struct{}) {
	data := 0

	for {
		select {
		case <-interrupt:
			counter.handle(0)
			return
		default:
			data++
			counter.handle(data)
			log.Println("Counter value:", data)
			time.Sleep(5 * time.Second)
		}
	}
}

func (counter *Counter) initDevice() {
	interrupt := make(chan struct{})

	for {
		select {
		case status := <-counter.status:
			if status == ON {
				go counter.runDevice(interrupt)
			}
			if status == OFF {
				interrupt <- struct{}{}
			}
		}
	}
}

func (counter *Counter) TurnOn() {
	counter.status <- ON
	log.Println("device on")
}

func (counter *Counter) TurnOff() {
	counter.status <- OFF
	log.Println("device off")
}

func NewCounter(h func (x int)) *Counter {
	counter := &Counter{
		status: make(chan int),
		handle: h,
	}

	go counter.initDevice()

	log.Println("New conter device")
	return counter
}

func CloseCounter(counter *Counter) {
	close(counter.status)
}
