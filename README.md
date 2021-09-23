# GoPigeon

GoPigeon is a MQTT broker implementation in Go.

## Features

GoPigeon implements the [MQTT 3.1.1.](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html) specification and its `mqtt` package provides some reusable solutions for encoding/decoding and control packet handling.

## Work In Progress

The implementation of the specification is tracked via issues. The project is in the earliest stages and is currently being developed. At the moment, the broker can keep a table of subscribers and publish messages.

Please see [mosquitto_pub](https://mosquitto.org/man/mosquitto_pub-1.html) and [mosquitto_sub](https://mosquitto.org/man/mosquitto_sub-1.html) for MQTT clients that can be used to test this broker.
