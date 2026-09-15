// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

//go:build android && !ts_omit_ssh

package main

import (
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSSHExecTargetTermux(t *testing.T) {
	sshExe := "/data/data/com.termux/files/usr/bin/ssh"
	argv := []string{sshExe, "-o", "ProxyCommand=tailcat address 22", "host"}
	gotExe, gotArgv := sshExecTarget(sshExe, argv)
	wantExe := "/system/bin/linker"
	if strconv.IntSize == 64 {
		wantExe += "64"
	}
	wantArgv := []string{sshExe, sshExe, "-o", "ProxyCommand=tailcat address 22", "host"}
	if gotExe != wantExe || !reflect.DeepEqual(gotArgv, wantArgv) {
		t.Fatalf("sshExecTarget() = (%q, %q); want (%q, %q)", gotExe, gotArgv, wantExe, wantArgv)
	}
}

func TestSSHExecutableFromTermuxEnv(t *testing.T) {
	want := "/data/data/com.termux/files/usr/bin/tailcat"
	t.Setenv("TERMUX_EXEC__PROC_SELF_EXE", want)
	got, err := sshExecutable()
	if err != nil || got != want {
		t.Fatalf("sshExecutable() = %q, %v; want %q", got, err, want)
	}
}

func TestTermuxWrappedExecutable(t *testing.T) {
	if os.Getenv("TAILCAT_TEST_WRAPPED") == "1" {
		exe, err := sshExecutable()
		if err != nil {
			t.Fatal(err)
		}
		if want := os.Getenv("TERMUX_EXEC__PROC_SELF_EXE"); exe != want {
			t.Fatalf("sshExecutable() = %q; want %q", exe, want)
		}
		return
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	linker := "/system/bin/linker"
	if strconv.IntSize == 64 {
		linker += "64"
	}
	cmd := exec.Command(linker, exe, "-test.v", "-test.run=^TestTermuxWrappedExecutable$")
	cmd.Args[0] = exe
	cmd.Env = append(os.Environ(), "TERMUX_VERSION=test", "TAILCAT_TEST_WRAPPED=1", "TERMUX_EXEC__PROC_SELF_EXE=/data/data/com.termux/files/usr/bin/tailcat")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("wrapped test failed: %v: %s", err, out)
	}
	if !strings.Contains(string(out), "--- PASS: TestTermuxWrappedExecutable") {
		t.Fatalf("wrapped test did not run: %s", out)
	}
}

func TestSSHExecTargetOutsideTermux(t *testing.T) {
	sshExe := "/system/bin/ssh"
	argv := []string{sshExe, "-V"}
	gotExe, gotArgv := sshExecTarget(sshExe, argv)
	if gotExe != sshExe || !reflect.DeepEqual(gotArgv, argv) {
		t.Fatalf("sshExecTarget() = (%q, %q); want (%q, %q)", gotExe, gotArgv, sshExe, argv)
	}
}
