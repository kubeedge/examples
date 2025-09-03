package driver

import (
	"sync"
	"time"

	"k8s.io/klog/v2"
)

const (
	ON = iota
	OFF
	PAUSED
	
)

type Counter struct {
	statusChan chan int
	handle func (int)
	curValue int
	curStatus int
    mu sync.Mutex
	doneChan chan struct{}
	wg sync.WaitGroup
}


func (counter *Counter) runDevice() {
	defer counter.wg.Done()
	ticker :=time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-counter.doneChan:
			return
		case status:=<-counter.statusChan:
			counter.mu.Lock()
			switch status {
			case ON:
				if counter.curStatus==OFF {
					counter.curValue=0
				}
				
				counter.curStatus=ON
			case OFF:
				counter.curValue=0
				counter.curStatus=OFF
			case PAUSED:
				counter.curStatus=PAUSED
			}
			counter.mu.Unlock()
		case <-ticker.C:
			counter.mu.Lock()
			if counter.curStatus == ON {
				counter.curValue++
				counter.handle(counter.curValue)
			}
			counter.mu.Unlock()
		}
	}

}


func (counter *Counter) ResetValue(){
	counter.mu.Lock()
	defer counter.mu.Unlock()
	counter.curValue=0
}

func (counter *Counter) SetCurValue(x int)  {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	counter.curValue=x
}
func (counter *Counter) GetCurValue() int {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	return counter.curValue
}

func (counter *Counter) GetStatus() int {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	return counter.curStatus
}

func (counter *Counter) SetStatus(status int) {
	counter.mu.Lock()
	defer counter.mu.Unlock()
	if status!=counter.curStatus {
		counter.statusChan<-status
	}
}



func NewCounter(h func (x int)) *Counter {
	counter := &Counter{
		statusChan: make(chan int),
		handle: h,
		curStatus: OFF,
		curValue: 0,
		mu: sync.Mutex{},
		doneChan: make(chan struct{}),
		wg: sync.WaitGroup{},
	}
	klog.Infoln("Counter Device start running")
	counter.wg.Add(1)
	go counter.runDevice()
	
	counter.statusChan<-counter.curStatus
	return counter
}

func CloseCounter(counter *Counter) {
	//provide a signal to stop the device
	close(counter.doneChan)

	//wait for the goroutine return when statusChan stopping all operations
	counter.wg.Wait()

	close(counter.statusChan)

	counter.mu.Lock()
	counter.curValue=0
	counter.curStatus=OFF
	counter.mu.Unlock()

	klog.Infoln("Counter resources released completely")
}
