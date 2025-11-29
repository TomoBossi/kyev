package main

import (
	"flag"
	"fmt"
)

type flags struct {
	nameMatch string
	physMatch string
}

func newFlags() (*flags, error) {
	nameMatch := flag.String("name", "", "Narrow the keyboard search to devices with names containing this substring, ignoring capitalization. It is REQUIRED. Device names are exactly as listed by xinput --list.")
	physMatch := flag.String("phys", "", "Narrow the keyboard search to devices with physical interface data containing this substring, ignoring capitalization. For example, \"usb\".")
	flag.Parse()

	if *nameMatch == "" {
		return nil, fmt.Errorf("name match substring not provided")
	}

	return &flags{
		nameMatch: *nameMatch,
		physMatch: *physMatch,
	}, nil
}
