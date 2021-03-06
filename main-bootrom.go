// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// +build bootrom

// First stage bootstrap. Looks for second stage or recovers the system
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/platinasystems/goes/external/partitions"
)

func copyFile(src, dst string, perm os.FileMode) (err error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("Error reading %s: %s", src, err)
	}
	err = ioutil.WriteFile(dst, b, perm)
	if err != nil {
		return fmt.Errorf("Error writing %s: %s", dst, err)
	}
	return nil
}

const name = "goes-bootrom"

func setupGoesBoot(dev string) (err error) {
	err = syscall.Mount("none", "/dev", "devtmpfs", 0, "")
	if err != nil {
		fmt.Printf("syscall.Mount(/dev) failed: %s\n", err)
		return
	}

	defer func() {
		err = syscall.Unmount("/dev", 0)
		if err != nil {
			fmt.Printf("syscall.Unmount(/dev) failed: %s\n", err)
		}
	}()

	sb, err := partitions.ReadSuperBlock(dev)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("partitions.ReadSuperBlock(%s) failed: %s\n",
				dev, err)
		}
		return
	}
	if sb == nil {
		fmt.Printf("Unable to recognize partition format on %s\n",
			dev)
		return
	}
	err = syscall.Mount(dev, "/boot", sb.Kind(), syscall.MS_RDONLY, "")
	if err != nil {
		fmt.Printf("syscall.Mount(%s[%s]) failed: %s\n", dev, sb.Kind(),
			err)
		return
	}

	defer func() {
		err = syscall.Unmount("/boot", 0)
		if err != nil {
			fmt.Printf("syscall.Unmount(%s) failed: %s\n",
				dev, err)
			return
		}
	}()

	path := "/boot/boot/" + packageName
	if _, err := os.Stat(path + "/" + packageName); err != nil {
		path = "/boot/" + packageName
		if _, err := os.Stat(path + "/" + packageName); err != nil {
			return nil
		}
	}

	for _, c := range []struct {
		src  string
		dst  string
		perm os.FileMode
	}{
		{path + "/" + packageName, "/sbin/goes-boot", 0755},
		{path + "/init", "/etc/goes/init", 0644},
		{path + "/start", "/etc/goes/start", 0644},
		{path + "/stop", "/etc/goes/stop", 0644},
		{path + "/authorized_keys", "/etc/goes/sshd/authorized_keys", 0600},
		{path + "/resolv.conf", "/etc/resolv.conf", 0644},
	} {
		err = copyFile(c.src, c.dst, c.perm)
		if err != nil {
			fmt.Printf("Error in copyFile: %s\n", err)
			return
		}
	}
	return nil
}

func execGoesBoot(dev string) {
	if err := setupGoesBoot(dev); err != nil {
		return
	}
	if _, err := os.Stat("/sbin/goes-boot"); err != nil {
		return
	}

	fmt.Printf("trying goes-boot (%s): args %v\n", dev, os.Args)
	err := syscall.Exec("/sbin/goes-boot", os.Args, os.Environ())
	fmt.Printf("syscall.Exec (%s) failed: %s\n", dev, err)
	return
}

func main() {
	if os.Args[0] == "/init" {
		execGoesBoot("/dev/sdb1")
		execGoesBoot("/dev/sda6")
		execGoesBoot("/dev/sda1")
		execGoesBoot("/dev/sda2")
	}

	args := os.Args
	if filepath.Base(args[0]) == name {
		args[0] = name
	}

	if err := Goes.Main(os.Args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
