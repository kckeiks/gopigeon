package gopigeon

func HandlePingreq(c *Client, fh *FixedHeader) error {
	if fh.Flags != 0 {
		return PingreqReservedFlagError
	}
	pingresp := []byte{Pingres << 4, 0}
	_, err := c.Conn.Write(pingresp)
	if err != nil {
		return err
	}
	return nil
}
