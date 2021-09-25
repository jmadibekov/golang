package main

import (
	"io"
)

type rot13Reader struct {
	r io.Reader
}

func rot13(b byte) byte {
	if b >= 65 && b <= 90 {
		b += 13
		if b > 90 {
			b = 65 + b - 91
		}
	} else if b >= 97 && b <= 122 {
		b += 13
		if b > 122 {
			b = 97 + b - 123
		}
	}
	return b
}
func (reader *rot13Reader) Read(b []byte) (n int, err error) {
	n, err = reader.r.Read(b)

	//	fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
	//	fmt.Printf("b[:n] = %q\n", b[:n])

	if err == io.EOF {
		return 0, err
	}

	for i := 0; i < n; i++ {
		b[i] = rot13(b[i])
	}

	return n, err
}
