package brewgadgets


import (
	"bitbucket.com/cswank/gogadgets"
	"log"
)


type HLT struct {
	gogadgets.InputDevice
	GPIO gogadgets.Poller
	Value float64
	Units string
	out chan<- gogadgets.Value
}

func NewHLT(pin *gogadgets.Pin) (gogadgets.InputDevice, error) {
	var err error
	var h *HLT
	pin.Edge = "rising"
	gpio, err := gogadgets.NewGPIO(pin)
	if err == nil {
		h = &HLT{
			GPIO:gpio.(gogadgets.Poller),
			Value: pin.Value.(float64),
			Units: pin.Units,
		}
	}
	return h, err
}

func (h *HLT) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Value) {
	h.out = out
	value := make(chan float64)
	err := make(chan error)
	keepGoing := true
	for keepGoing {
		go h.wait(value, err)
		select {
		case msg := <- in:
			keepGoing = h.ReadMessage(msg)
		case val := <-value:
			h.Value = val
			h.SendValue()
		case e := <-err:
			log.Println(e)
		}
	}
}

func (h *HLT) wait(out chan<- float64, err chan<- error) {
	val, e := h.GPIO.Wait()
	if e != nil {
		err<- e
	} else {
		if val {
			out<- h.Value
		} else {
			out<- 0.0
		}
	}
}

func (h *HLT) ReadMessage(msg gogadgets.Message) (keepGoing bool) {
	keepGoing = true
	if msg.Type == "command" && msg.Body == "shutdown" {
		keepGoing = false
	} else if msg.Type == "command" && msg.Body == "status" {
		h.SendValue()		
		
	} else if msg.Type == "command" && msg.Body == "hlt volume change" {
		h.Value += msg.Value.Value.(float64)
		h.SendValue()
	}
	return keepGoing
}

func (h *HLT) SendValue() {
	h.out<- gogadgets.Value{
		Value: h.Value,
		Units: h.Units,
	}
}
