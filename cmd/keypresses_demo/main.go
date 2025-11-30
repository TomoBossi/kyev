package main

import (
	"flag"
	"fmt"

	"github.com/tomobossi/kyev"
)

func main() {
	flags, err := newFlags()
	if err != nil {
		flag.Usage()
		fmt.Println()
		panic(err)
	}

	keyboard, err := kyev.GetKeyboard(flags.nameMatch, flags.physMatch)
	if err != nil {
		panic(err)
	}

	fmt.Println("press ESC to exit")
	for {
		keypresses, err := keyboard.GetKeyPresses()
		if err != nil {
			panic(err)
		}

		escPressed := false
		for _, keypress := range keypresses {
			fmt.Printf("Code: %d\nTime: %d.%d\n\n", keypress.Code, keypress.Time.Sec, keypress.Time.Usec)
			if keypress.Code == kyev.KEY_ESC {
				escPressed = true
				break
			}
		}

		if escPressed {
			break
		}
	}
}
