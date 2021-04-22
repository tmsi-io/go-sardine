package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// GetLocalIPList
// Get local host ip list without loopback and multicast
func GetLocalIPList() (IPList, error) {
	var ips IPList
	nicS, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("Get Address: %+v\n", err.Error()))
		return IPList{}, err
	}
	for _, i := range nicS {
		if strings.HasPrefix(i.Name, "lo") { // ignore local nic
			fmt.Printf("ignore nic [%s] \n", i.Name)
			continue
		}
		ipAddr, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("local Addresses: %+v\n", err.Error()))
			continue
		}
		for _, ip := range ipAddr {
			// ignore loop back and multicast ip addrs
			if Obj, ok := ip.(*net.IPNet); ok && !Obj.IP.IsLoopback() && !Obj.IP.IsMulticast() {
				if Obj.IP.To4() != nil {
					ips.IPv4 = append(ips.IPv4, Obj.IP.String())
				}
			}
		}
	}
	return ips, nil
}

// GetLocalIPByNicName
// filter nic ip for spic nic name
func GetLocalIPByNicName(nicName string) (string, error) {
	phyNic, err := net.InterfaceByName(nicName)
	if err != nil {
		log.Println(err)
		return "", err
	}
	addrS, _ := phyNic.Addrs()
	for _, ip := range addrS {
		if Obj, ok := ip.(*net.IPNet); ok && !Obj.IP.IsLoopback() && !Obj.IP.IsMulticast() {
			if Obj.IP.To4() != nil {
				return Obj.IP.String(), nil
			}
		}
	}
	return "", errors.New("No Match Nic Find. ")
}

// GetMACByNicName
// Get MAC address by nic name
func GetMACByNicName(nicName string) (string, error) {
	phyNic, err := net.InterfaceByName(nicName)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("%s", phyNic.HardwareAddr), nil
}

// GetMetricNicIP
// get default route nic ip
func GetMetricNicIP() (string, error) {
	return "", nil
}
