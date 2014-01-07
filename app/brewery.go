package main


import (
	"bitbucket.com/cswank/gogadgets"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/brewery/brewgadgets"
	"encoding/json"
	"io/ioutil"
	"flag"
)

var (
	configFlag = flag.String("c", "", "Path to the config json file")
)

func main() {
	flag.Parse()
	if !utils.FileExists("/sys/bus/w1/devices/28-0000047ade8f") {
		ioutil.WriteFile("/sys/devices/bone_capemgr.9/slots", []byte("BB-W1:00A0"), 0666)
	}
	b, err := ioutil.ReadFile(*configFlag)
	if err != nil {
		panic(err)
	}
	cfg := &gogadgets.Config{}
	err = json.Unmarshal(b, cfg)
	a := gogadgets.NewApp(cfg)

	config := &brewgadgets.MashConfig{
		TankRadius: 7.5 * 2.54,
		ValveRadius: 0.1875 * 2.54,
		Coefficient: 0.43244,
	}
	mashVolume, _ := brewgadgets.NewMash(config)

	mash := &gogadgets.Gadget{
		Location: "mash tun",
		Name: "volume",
		Input: mashVolume,
		Direction: "input",
		OnCommand: "n/a",
		OffCommand: "n/a",
		UID: "mash volume",
	}
	
	a.AddGadget(mash)
	poller, err := gogadgets.NewGPIO(&gogadgets.Pin{Port:"8", Pin:"9", Direction:"in", Edge:"rising"})
	if err != nil {
		panic(err)
	}
	hltVolume := &brewgadgets.HLT{
		GPIO: poller.(gogadgets.Poller),
		Value: 26.5,
		Units: "liters",
	}

	hlt := &gogadgets.Gadget{
		Location: "hlt",
		Name: "volume",
		Input: hltVolume,
		Direction: "input",
		OnCommand: "n/a",
		OffCommand: "n/a",
		UID: "hlt volume",
	}
	a.AddGadget(hlt)
	stop := make(chan bool)
	a.Start(stop)
}
