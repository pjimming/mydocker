package network

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type BridgeNetworkDriver struct {
}

func (d *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	n := &Network{
		Name:    name,
		IpRange: ipRange,
		Driver:  d.Name(),
	}

	if err := d.initBridge(n); err != nil {
		logrus.Errorf("[Create] init %s bridge error, %v", name, err)
		return nil, err
	}
	return n, nil
}

func (d *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (d *BridgeNetworkDriver) Delete(network Network) error {
	// 根据名字找到对应的Bridge设备
	br, err := netlink.LinkByName(network.Name)
	if err != nil {
		return err
	}
	// 删除网络对应的 Lin ux Bridge 设备
	return netlink.LinkDel(br)
}

// Connect 连接一个网络和网络端点
func (d *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	bridgeName := network.Name
	// 通过接口名获取到 Linux Bridge 接口的对象和接口属性
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	// 创建 Veth 接口的配置
	la := netlink.NewLinkAttrs()
	// 由于 Linux 接口名的限制,取 endpointID 的前
	la.Name = endpoint.ID[:5]
	// 通过设置 Veth 接口 master 属性，设置这个Veth的一端挂载到网络对应的 Linux Bridge
	la.MasterIndex = br.Attrs().Index
	// 创建 Veth 对象，通过 PeerNarne 配置 Veth 另外 端的接口名
	// 配置 Veth 另外 端的名字 cif {endpoint ID 的前 位｝
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + endpoint.ID[:5],
	}
	// 调用netlink的LinkAdd方法创建出这个Veth接口
	// 因为上面指定了link的MasterIndex是网络对应的Linux Bridge
	// 所以Veth的一端就已经挂载到了网络对应的LinuxBridge.上
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("error Add Endpoint Device: %v", err)
	}
	// 调用netlink的LinkSetUp方法，设置Veth启动
	// 相当于ip link set xxx up命令
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("error Add Endpoint Device: %v", err)
	}
	return nil
}

func (d *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}

// initBridge 初始化Linux Bridge
/*
Linux Bridge 初始化流程如下：
* 1）创建 Bridge 虚拟设备
* 2）设置 Bridge 设备地址和路由
* 3）启动 Bridge 设备
* 4）设置 iptables SNAT 规则
*/
func (d *BridgeNetworkDriver) initBridge(n *Network) error {
	bridgeName := n.Name
	// 1）创建 Bridge 虚拟设备
	if err := createBridgeInterface(bridgeName); err != nil {
		return err
	}

	// 2）设置 Bridge 设备地址和路由
	gatewayIP := *n.IpRange
	gatewayIP.IP = n.IpRange.IP

	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return err
	}
	// 3）启动 Bridge 设备
	if err := setInterfaceUP(bridgeName); err != nil {
		return err
	}

	// 4）设置 iptables SNAT 规则
	if err := setupIPTables(bridgeName, n.IpRange); err != nil {
		return err
	}

	return nil
}

// createBridgeInterface 创建Bridge设备
// ip link add xxxx
func createBridgeInterface(bridgeName string) error {
	// 先检查是否己经存在了这个同名的Bridge设备
	if _, err := net.InterfaceByName(bridgeName); err == nil || !errors.Is(err, errNoSuchInterface) {
		logrus.Errorf("[createBridgeInterface] net.InterfaceByName error, %v", err)
		return err
	}

	// create *netlink.Bridge object
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName
	// 使用刚才创建的Link的属性创netlink Bridge对象
	br := &netlink.Bridge{LinkAttrs: la}
	// 调用 net link Linkadd 方法，创 Bridge 虚拟网络设备
	// netlink.LinkAdd 方法是用来创建虚拟网络设备的，相当于 ip link add xxxx
	if err := netlink.LinkAdd(br); err != nil {
		logrus.Errorf("[createBridgeInterface] LinkAdd error, %v", err)
		return err
	}
	return nil
}

// Set the IP addr of a netlink interface
// ip addr add xxx命令
func setInterfaceIP(name, rawIP string) error {
	var iface netlink.Link
	var err error
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		logrus.Infof("[setInterfaceIP] retrieving new bridge netlink link [%s], %v", name, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logrus.Errorf("[setInterfaceIP] link [%s] bridge error, %v", name, err)
		return err
	}
	// 由于 netlink.ParseIPNet 是对 net.ParseCIDR一个封装，因此可以将 net.PareCIDR中返回的IP进行整合
	// 返回值中的 ipNet 既包含了网段的信息，192 168.0.0/24 ，也包含了原始的IP 192.168.0.1
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		logrus.Errorf("[setInterfaceIP] parse ip[%s] net error, %v", rawIP, err)
		return err
	}
	// 通过  netlink.AddrAdd给网络接口配置地址，相当于ip addr add xxx命令
	// 同时如果配置了地址所在网段的信息，例如 192.168.0.0/24
	// 还会配置路由表 192.168.0.0/24 转发到这 testbridge 的网络接口上
	addr := &netlink.Addr{IPNet: ipNet}
	return netlink.AddrAdd(iface, addr)
}

// setInterfaceUP 启动Bridge设备
// 等价于 ip link set xxx up 命令
func setInterfaceUP(interfaceName string) error {
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}
	// 等价于 ip link set xxx up 命令
	if err = netlink.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}

// setupIPTables 设置 iptables 对应 bridge MASQUERADE 规则
// iptables -t nat -A POSTROUTING -s 172.18.0.0/24 -o eth0 -j MASQUERADE
// iptables -t nat -A POSTROUTING -s {subnet} -o {deviceName} -j MASQUERADE
func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	// 拼接命令
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	// 执行该命令
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("iptables Output, %v", output)
		return err
	}
	return nil
}
