package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func findIPInNetworkFromIface(dstIP net.IP, iface netlink.Link) (net.IP, error) {
	addrs, err := netlink.AddrList(iface, netlink.FAMILY_V4)
	if err != nil {
		return nil, err
	}
	for _, a := range addrs {
		if a.Contains(dstIP) {
			return a.IP, nil
		}

	}
	/*
		addrs, err := iface.Addrs()

		if err != nil {
			return nil, err
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ipnet.Contains(dstIP) {
					return ipnet.IP, nil
				}
			}
		}
	*/
	return nil, fmt.Errorf("iface: '%s' can't reach ip: '%s'", iface.Attrs().Name, dstIP)
}

func findUsableInterfaceForNetwork(dstIP net.IP) (net.IP, *netlink.Link, error) {
	/*
		ifaces, err := net.Interfaces()
		if err != nil {
			return nil, err
		}
	*/
	//find interfaces
	ifaces, err := netlink.LinkList()
	if err != nil {
		return nil, nil, err
	}

	isDown := func(iface netlink.Link) bool {
		return iface.Attrs().Flags&1 == 0
	}

	verboseLog.Println("search usable interface")
	logIfaceResult := func(msg string, iface netlink.Link) {
		verboseLog.Printf("%10s: %6s %18s  %s", msg, iface.Attrs().Name, iface.Attrs().HardwareAddr, iface.Attrs().Flags)
	}

	for _, iface := range ifaces {
		if isDown(iface) {
			logIfaceResult("DOWN", iface)
			continue
		}

		if srcIP, err := findIPInNetworkFromIface(dstIP, iface); err != nil {
			logIfaceResult("OTHER NET", iface)
			continue
		} else {

			logIfaceResult("USABLE", iface)
			return srcIP, &iface, nil
		}
	}

	return nil, nil, errors.New("no usable interface found")
}
