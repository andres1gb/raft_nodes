package raftnodes

import (
	"errors"
	"fmt"
	"net/netip"
)

type Address struct {
	addr netip.Addr
	port uint64
}

func NewAddress(addr string, port uint64) (*Address, error) {
	address, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	if port == 0 || port > 65535 {
		return nil, errors.New("port number should be in the 1:65535 range")
	}
	return &Address{
		addr: address,
		port: port,
	}, nil
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.addr.String(), a.port)
}
