package lib

type ConnectPkt struct {
    pktType byte 
    protocolName []byte
    protocolLevel byte
    userNameFlag byte
    psswdFlag byte
    willRetainFlag byte
    willQoSFlag byte
    willFlag byte
    cleanSession byte
    keepAlive []byte    
}

func (cp ConnectPkt) decode() {
	
}