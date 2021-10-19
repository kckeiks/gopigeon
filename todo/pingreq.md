# 3.12 PINGREQ – PING request

## `THIS IS DONE`

The PINGREQ Packet is sent from a Client to the Server. It can be used to:

Indicate to the Server that the Client is alive in the absence of any other Control Packets being sent from the Client to the Server.
Request that the Server responds to confirm that it is alive.
Exercise the network to indicate that the Network Connection is active.

This Packet is used in Keep Alive processing, see Section 3.1.2.10 for more details.

## 3.12.1 Fixed header

Figure 3.33 – 
PINGREQ Packet fixed header
 
## 3.12.2 Variable header

The PINGREQ Packet has no variable header.
3.12.3 Payload

The PINGREQ Packet has no payload.
3.12.4 Response

`The Server MUST send a PINGRESP Packet in response to a PINGREQ Packet [MQTT-3.12.4-1].`
