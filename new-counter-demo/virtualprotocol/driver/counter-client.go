package driver


const (
	BUFFER=5
)

type CounterClient interface {
	InitCounterCli()   
	GetCounter()int 
	SetCurValue(int)  
	GetStatus()int  
	SetStatus(int) 
	ResetValue()
	Close() 		
}

type Client struct{
	counter *Counter
	dataChan chan int
}

func (c *Client) InitCounterCli() {
	c.dataChan = make(chan int,BUFFER)
	// func to handle the data from counter (also can change to mqtt)
	c.counter= NewCounter(func(x int) {
		if c.dataChan!=nil{
			select {
			case c.dataChan <- x:
			default:
				// when the chan is full, drop the old data and send new data 
				<-c.dataChan
				c.dataChan <- x
			}
		}

	})
}

func (c *Client) ResetValue() {
	c.counter.ResetValue()
}

func (c *Client) GetCounter() int {
	select {
	case x := <-c.dataChan:
		return x
	default:
		return c.counter.GetCurValue()
	}
}
func (c *Client) SetCurValue(x int) {
	c.counter.SetCurValue(x)
}

func (c *Client) GetStatus() int {
	return c.counter.GetStatus()
}
func (c *Client) SetStatus(status int) {
	c.counter.SetStatus(status)
}
func (c *Client) Close() {
	CloseCounter(c.counter)
}

func NewCounterClient() *Client {
	return &Client{
	}
}

