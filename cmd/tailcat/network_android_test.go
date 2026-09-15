// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build android

package main

import (
	"net"
	"os"
	"testing"

	"tailscale.com/net/netmon"
	"tailscale.com/util/eventbus"
)

func TestTermuxNetmon(t *testing.T) {
	if os.Getenv("TERMUX_VERSION") == "" {
		t.Skip("requires Termux network restrictions")
	}

	monitor, err := netmon.New(eventbus.New(), func(string, ...any) {})
	if err != nil {
		t.Fatalf("netmon.New() in Termux: %v", err)
	}
	defer monitor.Close()
	if monitor.InterfaceState() == nil {
		t.Fatal("netmon.New() returned a nil interface state")
	}
}

func TestTermuxInterfacesFromIPs(t *testing.T) {
	if got := termuxInterfacesFromIPs(nil, nil); len(got) != 0 {
		t.Fatalf("no route: got %d interfaces; want none", len(got))
	}
	got := termuxInterfacesFromIPs(net.IPv4(192, 0, 2, 10), net.ParseIP("2001:db8::10"))
	if len(got) != 1 {
		t.Fatalf("got %d interfaces; want one", len(got))
	}
	if !got[0].IsUp() || len(got[0].AltAddrs) != 2 {
		t.Fatalf("interface = %+v; want up with two addresses", got[0])
	}
}
