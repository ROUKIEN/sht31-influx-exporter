package sht31

import (
	"testing"
)

func TestToTemperatureCelsius(t *testing.T) {
	bytes := []byte{96, 37, 167, 153, 152, 143}

	result := ToTemperatureCelsius(bytes)
	const expected = 20.724800
	if result != expected {
		t.Errorf("It failed. Expected %f, got %f", expected, result)
	}
}

func TestToRelativeHumidity(t *testing.T) {
	bytes := []byte{96, 37, 167, 153, 152, 143}

	result := ToRelativeHumidity(bytes)
	const expected = 59.998474
	if result != expected {
		t.Errorf("It failed. Expected %f, got %f", expected, result)
	}
}
