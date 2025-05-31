package wss

import (
	"net"
)

type DirectClient struct {
	addr string // target address
}

func (client *DirectClient) ProxyType() int {
	return ProxyTypeDirect
}

func (client *DirectClient) Trigger(data []byte) bool {
	return true // always true, direct connection
}
func (client *DirectClient) EstablishData(origin []byte) ([]byte, error) {
	return nil, nil // no data to send in establishing step
}
func (client *DirectClient) ParseHeader(conn net.Conn, header []byte) (string, error) {
	// direct connection, no header parsing needed
	return client.addr, nil
}

var _ ProxyInterface = &DirectClient{}
