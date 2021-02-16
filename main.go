package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"rpi-sensors/devices/sht31"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/pin"
	"periph.io/x/periph/conn/pin/pinreg"
	"periph.io/x/periph/host"
)

// InfluxToken is injected at build time.
var InfluxToken string

// InfluxEndpoint is injected at build time.
var InfluxEndpoint string

// InfluxOrg is injected at build time.
var InfluxOrg string

// InfluxBucket is injected at build time.
var InfluxBucket string

const deviceID = "office"

func printPin(fn string, p pin.Pin) {
	name, pos := pinreg.Position(p)
	if name != "" {
		log.Printf(" %-4s: %-10s found on header %s, #%d\n", fn, p, name, pos)
	} else {
		log.Printf(" %-4s: %-10s\n", fn, p)
	}
}

func main() {
	influxClient := influxdb2.NewClient(InfluxEndpoint, InfluxToken)
	defer influxClient.Close()
	writeAPI := influxClient.WriteAPI(InfluxOrg, InfluxBucket)

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

	deviceConn := &i2c.Dev{Addr: sht31.MainI2CAddress, Bus: bus}
	// sht31.NewI2C(bus)

	stopPeriodicMeasureCmd := []byte{0x30, 0x93}
	if err := deviceConn.Tx(stopPeriodicMeasureCmd, nil); err != nil {
		log.Fatal(err)
	}

	startPeriodicMeasureCmd := []byte{sht31.CmdPeriodicMeasurementOnePerSecMsb, 0x30}
	if err := deviceConn.Tx(startPeriodicMeasureCmd, nil); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	i := 0
	t := time.NewTicker(1 * time.Second)
	for l := gpio.Low; ; l = !l {
		readoutPeriodicMeasureCmd := []byte{sht31.CmdPeriodicReadoutMsb, sht31.CmdPeriodicReadoutLsb}
		values := make([]byte, 6) // First two bytes byte are temp, Third are temp CRC, Fourth & Fivest bytes are humidity & Sixth is humidity CRC
		if err := deviceConn.Tx(readoutPeriodicMeasureCmd, values); err != nil {
			fmt.Println("oh no !")
			log.Fatal(err)
		}
		// @TODO ensure CRC checksums are correct
		// fmt.Println(fmt.Sprintf("Temperature: %f°C; Humidity: %f%%", sht31.ToTemperatureCelsius(values), sht31.ToRelativeHumidity(values)))

		// write line protocol
		writeAPI.WriteRecord(fmt.Sprintf("temperature,device=%s value=%f", deviceID, sht31.ToTemperatureCelsius(values)))
		writeAPI.WriteRecord(fmt.Sprintf("humidity,device=%s value=%f", deviceID, sht31.ToRelativeHumidity(values)))

		// Flush writes
		if i > 5*60 {
			writeAPI.Flush()
			i = 0
		} else {
			i++
		}
		<-t.C
	}
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// rpi.P1_40.Out(gpio.Low)
		// Run Cleanup
		fmt.Println("Gracefully stopping the app")
		os.Exit(1)
	}()
}
