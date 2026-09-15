// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build !ts_omit_ssh

package main

import (
	"os"
	"runtime"
	"strings"
)

// sshExecutable returns tailcat's path for OpenSSH's ProxyCommand. On Android,
// Termux runs tailcat through the system linker, so os.Executable reports the
// linker instead; termux-exec records the original path in the environment.
func sshExecutable() (string, error) {
	if runtime.GOOS == "android" {
		if exe := os.Getenv("TERMUX_EXEC__PROC_SELF_EXE"); strings.HasPrefix(exe, "/data/data/com.termux/files/") {
			return exe, nil
		}
	}
	return os.Executable()
}
