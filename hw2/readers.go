// https://tour.golang.org/methods/22

package main

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.
func (reader MyReader) Read(b []byte) (n int, err error) {
	b[0] = 'A'
	return 1, nil
}
