package utils

import (
	"fmt"
	"net/netip"
)

type Address struct {
	Addr netip.Addr
	Port uint64
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.Addr.String(), a.Port)
}
