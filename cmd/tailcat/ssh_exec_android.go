// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build android && !ts_omit_ssh

package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

// execSSH replaces tailcat with the system SSH client. Android forbids direct
// exec of files in Termux's private data directory, so launch those binaries
// through the system dynamic linker, as termux-exec does for shell commands.
func execSSH(sshExe string, argv []string) error {
	exe, args := sshExecTarget(sshExe, argv)
	if exe != sshExe {
		if err := os.Setenv("TERMUX_EXEC__PROC_SELF_EXE", sshExe); err != nil {
			return err
		}
	}
	return syscall.Exec(exe, args, os.Environ())
}

func sshExecTarget(sshExe string, argv []string) (string, []string) {
	if !strings.HasPrefix(sshExe, termuxPrefix) {
		return sshExe, argv
	}
	linker := "/system/bin/linker"
	if strconv.IntSize == 64 {
		linker += "64"
	}
	args := make([]string, 0, len(argv)+1)
	args = append(args, argv[0], sshExe)
	args = append(args, argv[1:]...)
	return linker, args
}
