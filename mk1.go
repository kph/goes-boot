// Copyright © 2020 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// +build mk1

package main

import (
	"fmt"

	"github.com/platinasystems/ioport"
)

const packageName = "goes-boot-platina-mk1"
const recoveryUrl = "https://platina.io/goes/goes-boot-platina-mk1.cpio.xz"

func disableBootdog() (err error) {
	b, err := ioport.Inb(0x604)
	if err != nil {
		return fmt.Errorf("Error in Inb(0x604): %s", err)
	}
	b = b & 0xfd
	err = ioport.Outb(0x604, b)
	if err != nil {
		return fmt.Errorf("Error in Outb(0x604, %x): %s", b, err)
	}
	qspi := 0
	if b&0x80 != 0 {
		qspi = 1
	}
	fmt.Printf("Booted from QSPI%d\n", qspi)

	return
}
