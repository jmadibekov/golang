package main

import (
	"fmt"

	"io"
	"os"
	"strings"

	"golang.org/x/tour/reader"
)

func main() {
	// stringers
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	// errors
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))

	// readers
	reader.Validate(MyReader{})

	// rot13Reader
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
