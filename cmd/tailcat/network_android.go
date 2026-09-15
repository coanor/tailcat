// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build android

package main

import (
	"net"
	"os"

	"tailscale.com/net/netmon"
)

func init() {
	if os.Getenv("TERMUX_VERSION") != "" {
		netmon.RegisterInterfaceGetter(termuxInterfaces)
	}
}

func termuxInterfaces() ([]netmon.Interface, error) {
	return termuxInterfacesFromIPs(
		termuxLocalIP("udp4", "1.1.1.1:53"),
		termuxLocalIP("udp6", "[2606:4700:4700::1111]:53"),
	), nil
}

// A UDP connect only asks the kernel to select a local address; it does not
// send a packet. Android blocks the netlink query used by net.Interfaces for
// unprivileged Termux processes, but ordinary socket routing still works.
func termuxLocalIP(network, target string) net.IP {
	conn, err := net.Dial(network, target)
	if err != nil {
		return nil
	}
	defer conn.Close()
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil
	}
	return addr.IP
}

func termuxInterfacesFromIPs(ips ...net.IP) []netmon.Interface {
	addrs := make([]net.Addr, 0, len(ips))
	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			addrs = append(addrs, &net.IPNet{IP: ip4, Mask: net.CIDRMask(32, 32)})
		} else if ip6 := ip.To16(); ip6 != nil {
			addrs = append(addrs, &net.IPNet{IP: ip6, Mask: net.CIDRMask(128, 128)})
		}
	}
	if len(addrs) == 0 {
		return nil
	}
	return []netmon.Interface{{
		Interface: &net.Interface{Name: "termux", Flags: net.FlagUp},
		AltAddrs:  addrs,
	}}
}
