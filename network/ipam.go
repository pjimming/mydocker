package network

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// IPAM 用于网络 IP 地址的分配和释放，包括容器的IP地址和网络网关的IP地址
type IPAM struct {
	SubnetAllocatorPath string             // 分配文件存放位置
	Subnets             *map[string]string // 网段和位图算法的数组 map, key 是网段， value 是分配的位图数组
}

// 初始化一个IPAM的对象，默认使用/var/run/mydocker/network/ipam/subnet.json作为分配信息存储位置
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnets = &map[string]string{}

	if err = ipam.load(); err != nil {
		logrus.Errorf("[Allocate] load subnets error, %v", err)
		return nil, err
	}

	_, subnet, _ = net.ParseCIDR(subnet.String())
	one, size := subnet.Mask.Size()
	logrus.Debugf("[Allocate] one = %d, size = %d, subnet = %s", one, size, subnet)
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		// 用“0”填满这个网段的配置，uint8(size - one)表示这个网段中有多少个可用地址
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}

	for i := range (*ipam.Subnets)[subnet.String()] {
		if (*ipam.Subnets)[subnet.String()][i] == '0' {
			// 设置这个为“0”的序号值为“1” 即标记这个IP已经分配过了
			ipAlloc := []byte((*ipam.Subnets)[subnet.String()])
			ipAlloc[i] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipAlloc)
			ip = subnet.IP
			/*
				还需要通过网段的IP与上面的偏移相加计算出分配的IP地址，由于IP地址是uint的一个数组，
				需要通过数组中的每一项加所需要的值，比如网段是172.16.0.0/12，数组序号是65555,
				那么在[172,16,0,0] 上依次加[uint8(65555 >> 24)、uint8(65555 >> 16)、
				uint8(65555 >> 8)、uint8(65555 >> 0)]， 即[0, 1, 0, 19]， 那么获得的IP就
				是172.17.0.19.
			*/
			for t := uint(4); t > 0; t-- {
				[]byte(ip)[4-t] += uint8(i >> ((t - 1) * 8))
			}
			ip[3]++
			break
		}
	}
	if err = ipam.dump(); err != nil {
		logrus.Errorf("[Allocate] ipam dump error, %v", err)
	}
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipAddr *net.IP) error {
	ipam.Subnets = &map[string]string{}
	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		logrus.Errorf("[Release] ipam load error, %v", err)
		return err
	}
	// 和分配一样的算法，反过来根据IP找到位图数组中的对应索引位置
	c := 0
	releaseIP := ipAddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}
	// 然后将对应位置0
	ipAlloc := []byte((*ipam.Subnets)[subnet.String()])
	ipAlloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipAlloc)

	// 最后调用dump将分配结果保存到文件中
	err = ipam.dump()
	if err != nil {
		logrus.Error("[Release] dump ipam error", err)
	}
	return nil
}

func (ipam *IPAM) load() error {
	// 检查是否存在文件，如果不存在说明之前未分配
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("[IPAM.load] os.Stat error, %v", err)
			return err
		}
		logrus.Infof("[IPAM.load] %s not find", ipam.SubnetAllocatorPath)
		return nil
	}

	// 读取文件，加载配置信息
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	if err != nil {
		logrus.Errorf("[IPAM.load] open %s error, %v", ipam.SubnetAllocatorPath, err)
		return err
	}
	defer func() {
		_ = subnetConfigFile.Close()
	}()

	subnetJson, err := io.ReadAll(subnetConfigFile)
	if err != nil {
		logrus.Errorf("[IPAM.load] read %s error, %v", ipam.SubnetAllocatorPath, err)
		return err
	}
	if err = json.Unmarshal(subnetJson, ipam.Subnets); err != nil {
		logrus.Errorf("[IPAM.load] json unmarshal error, %v", err)
		return err
	}
	return nil
}

func (ipam *IPAM) dump() error {
	if err := os.MkdirAll(path.Dir(ipam.SubnetAllocatorPath), 0644); err != nil {
		logrus.Errorf("[IPAM.dump] mkdir %s error, %v", ipam.SubnetAllocatorPath, err)
		return err
	}
	// 打开存储文件 O_TRUNC 表示如果存在则消空，O_CREATE表示如果不存在则创建，O_WRONLY只写模式打开文件
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Errorf("[IPAM.dump] open %s error, %v", ipam.SubnetAllocatorPath, err)
		return err
	}
	defer func() {
		_ = subnetConfigFile.Close()
	}()

	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		logrus.Errorf("[IPAM.dump] marshal error, %v", err)
		return err
	}
	if _, err = subnetConfigFile.Write(ipamConfigJson); err != nil {
		logrus.Errorf("[IPAM.dump] write %s error, %v", ipam.SubnetAllocatorPath, err)
		return err
	}
	return nil
}
