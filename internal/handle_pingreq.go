package internal

import "github.com/kckeiks/gopigeon/mqttlib"

func handlePingreq(c *Client, fh *mqttlib.FixedHeader) error {
	if fh.Flags != 0 {
		return mqttlib.PingreqReservedFlagError
	}
	pingresp := []byte{mqttlib.Pingres << 4, 0}
	_, err := c.Conn.Write(pingresp)
	if err != nil {
		return err
	}
	return nil
}
