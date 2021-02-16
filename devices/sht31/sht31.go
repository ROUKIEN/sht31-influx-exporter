package sht31

import (
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
)

// Dev is a device.
type Dev struct {
	recurrentMeasureEnabled bool
	c                       conn.Conn
}

type Measure struct {
	Temperature int
	Humidity    int
}

// NewI2C returns a new Dev object that communicates over I2C with an SHT31 sensor.
func NewI2C(bus i2c.Bus) (*Dev, error) {
	// &i2c.Dev{Addr: ADDRESS, Bus: bus}
	return nil, nil
}

// EnableRecurrentMeasure sends a periodical measure signal to the sht31.
func (d *Dev) EnableRecurrentMeasure() error {
	d.recurrentMeasureEnabled = true

	startRecurrentMeasureCmd := []byte{CmdPeriodicMeasurementOnePerSecMsb, 0x30}
	if err := d.c.Tx(startRecurrentMeasureCmd, nil); err != nil {
		return err
	}

	return nil
}

// DisableRecurrentMeasure sends a command to the sht31 to stop recurrent measures.
func (d *Dev) DisableRecurrentMeasure() error {
	d.recurrentMeasureEnabled = false

	stopRecurrentMeasureCmd := []byte{0x30, 0x93}
	if err := d.c.Tx(stopRecurrentMeasureCmd, nil); err != nil {
		return err
	}

	return nil
}

// ReadoutMeasure asks the sht31 for the current measure.
func (d *Dev) ReadoutMeasure() (chan Measure, error) {
	timeTick := time.NewTicker(1 * time.Second)
	for {
		<-timeTick.C
	}
	return nil, nil
}

// ToTemperatureCelsius converts the byte payload into celsius temperature
func ToTemperatureCelsius(data []byte) float32 {
	temperature := float32(data[0])*256 + float32(data[1])
	return -45 + (175 * temperature / 65535.0)
}

// ToRelativeHumidity converts the byte payload into the relative humidity percentage
func ToRelativeHumidity(val []byte) float32 {
	return 100 * (float32(val[3])*256 + float32(val[4])) / 65535.0
}

const (
	MainI2CAddress = 0x44 // default I2C Address

	CmdMeasureClockStretchingEnabled  = 0x2c // measure with clock stretching
	CmdMeasureClockStretchingDisabled = 0x24 // measure with clock stretching

	CmdMeasureClockStretchingEnabledHighRepeatability    = 0x06 // measure with clock stretching & high repeatability
	CmdMeasureClockStretchingEnabledMediumRepeatability  = 0x0d // measure with clock stretching & medium repeatability
	CmdMeasureClockStretchingEnabledLowRepeatability     = 0x10 // measure with clock stretching & low repeatability
	CmdMeasureClockStretchingDisabledHighRepeatability   = 0x00 // measure without clock stretching & high repeatability
	CmdMeasureClockStretchingDisabledMediumRepeatability = 0x0b // measure without clock stretching & medium repeatability
	CmdMeasureClockStretchingDisabledLowRepeatability    = 0x16 // measure without clock stretching & low repeatability

	CmdMeasureHeater = 0x30 // Heater command prefix
	CmdHeaterEnable  = 0x6D // Enable heater
	CmdHeaterDisable = 0x66 // Disable heater

	CmdStatusRegisterMsb = 0xf32d // Status register MSB
	CmdStatusRegisterLsb = 0x2d   // Status register LSB

	CmdPeriodicMeasurementHalfPerSecMsb = 0x20
	CmdPeriodicMeasurementOnePerSecMsb  = 0x21
	CmdPeriodicMeasurementTwoPerSecMsb  = 0x22
	CmdPeriodicMeasurementFourPerSecMsb = 0x23
	CmdPeriodicMeasurementTenPerSecMsb  = 0x27

	CmdPeriodicReadoutMsb = 0xE0
	CmdPeriodicReadoutLsb = 0x00
)
