package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/pin"
	"periph.io/x/periph/conn/pin/pinreg"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

func printPin(fn string, p pin.Pin) {
	name, pos := pinreg.Position(p)
	if name != "" {
		log.Printf(" %-4s: %-10s found on header %s, #%d\n", fn, p, name, pos)
	} else {
		log.Printf(" %-4s: %-10s\n", fn, p)
	}
}

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	for _, ref := range i2creg.All() {
		fmt.Printf("- %s\n", ref.Name)
	}

	bus, err := i2creg.Open("")
	bus.SetSpeed(9600)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	if p, ok := bus.(i2c.Pins); ok {
		printPin("SCL", p.SCL())
		printPin("SDA", p.SDA())
	}

	device := &i2c.Dev{Addr: 0x44, Bus: bus}
	// readStatus := []byte{0xF32D}

	// is heater enabled ?
	/*
		if err := device.Tx(readStatus, nil); err != nil {
			log.Fatal(err)
		}
	*/

	t := time.NewTicker(1 * time.Second)
	for l := gpio.Low; ; l = !l {
		if err := rpi.P1_40.Out(l); err != nil {
			log.Fatal(err)
		}

		requestWrite := []byte{0x24}
		getValuesRead1 := make([]byte, 6)
		if err := device.Tx(requestWrite, getValuesRead1); err != nil {
			log.Fatal(err)
		}

		// fmt.Println("Sent values request")

		time.Sleep(200 * time.Millisecond)

		getValuesRead := make([]byte, 6)
		if err := device.Tx(nil, getValuesRead); err != nil {
			log.Fatal(err)
		}

		// fmt.Println(getValuesRead)
		<-t.C
	}
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		rpi.P1_40.Out(gpio.Low)
		// Run Cleanup
		os.Exit(1)
	}()
}
