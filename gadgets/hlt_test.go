package brewgadgets

import (
	"fmt"
	"testing"
	"time"
	"bitbucket.com/cswank/gogadgets"
)

type FakePoller struct {
	gogadgets.Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	if f.val {
		time.Sleep(10 * time.Second)
	} else {
		time.Sleep(100 * time.Millisecond)
		f.val = !f.val
	}
	return f.val, nil
}

func TestCreate(t *testing.T) {
	poller := &FakePoller{}
	_ = &HLT{
		GPIO: poller,
		volume: 5.0,
		units: "liters",
	}
}

func TestHLT(t *testing.T) {
	poller := &FakePoller{}
	h := &HLT{
		GPIO: poller,
		volume: 5.0,
		units: "liters",
	}
	out := make(chan gogadgets.Message)
	in := make(chan gogadgets.Value)
	go h.Start(out, in)
	val := <-in
	fmt.Println(val)
	if val.Value.(float64) != 5.0 {
		t.Error("should have been 5.0", val)
	}
	go func() {
		time.Sleep(10 * time.Millisecond)
		out<- gogadgets.Message{
			Type: "update",
			Body: "mash volume",
			Value: gogadgets.Value{
				Value: 0.5,
			},
		}
	}()
	val = <-in
	if val.Value.(float64) != 4.5 {
		t.Error("should have been 4.5", val)
	}
}
