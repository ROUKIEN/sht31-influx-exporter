# RPI SHT31

A golang program to read humidity/temperature from a SHT31 sensor (https://wiki.seeedstudio.com/Grove-TempAndHumi_Sensor-SHT31/) using periph.io (https://periph.io/) which exports data to an influxdb instance.

## Requirements

1. You'll need a raspberry pi 3 with I2C. The pi must be reachable over SSH.
2. Copy the .env.dist to .env and fill the mandatory values.

## Commands

### Build (ARM)

```shell
make build
```

The influxdb credentials are injected at build time through ldflags.

### Deploy

```shell
make deploy
```

> Once deployed, you may want to run the command as a service. Just create a new systemd unit from the systemd/sensor.service on the pi.

### Improvements

* I'd like to refactor the sht31 module to be fully compatible with periph.io.
* generate a custom linux kernel with buildroot/yocto to be able to minimize the footprint ?
