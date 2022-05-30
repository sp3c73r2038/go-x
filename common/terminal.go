package common

import (
	"fmt"

	"golang.org/x/term"
)

func AskPassword(prompt string) (rv string, err error) {
	fmt.Print(prompt)
	var pass []byte
	pass, err = term.ReadPassword(1)
	fmt.Println("")
	if err != nil {
		// cannot get key from console input, exit...
		return
	}
	rv = string(pass)
	return
}
