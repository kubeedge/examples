package counter

import (
	"fmt"
	"time"
)

const (
	ON = iota
	OFF
)

type Counter struct {
	Name          string
	CurrentStatus string
	status        chan int
	handle        func(string, int)
}

func (counter *Counter) runDevice(interrupt chan struct{}) {
	data := 0

	for {
		select {
		case <-interrupt:
			counter.handle(counter.Name, 0)
			return
		default:
			data++
			counter.handle(counter.Name, data)
			fmt.Printf("Counter: %s, Counter value: %d \n", counter.Name, data)
			time.Sleep(1 * time.Second)
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
}

func (counter *Counter) TurnOff() {
	counter.status <- OFF
}

func NewCounter(name string, h func(name string, x int)) *Counter {
	counter := &Counter{
		Name:          name,
		CurrentStatus: "OFF",
		status:        make(chan int),
		handle:        h,
	}

	go counter.initDevice()

	return counter
}

func CloseCounter(counter *Counter) {
	close(counter.status)
}
