package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskSize(t *testing.T) {
	_, subnet, _ := net.ParseCIDR("127.0.0.0/8")
	one, size := subnet.Mask.Size()
	t.Logf("one = %d, size = %d, subnet = %s", one, size, subnet.String())
}

func TestIPAM_Allocate(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.168.0.1/24")
	ip, err := ipAllocator.Allocate(ipNet)
	assert.Nil(t, err)
	t.Logf("allocate ip: %s", ip.String())
}
