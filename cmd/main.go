package main

import (
	"bytes"
	"fmt"

	"github.com/itchio/lzma"
)

func main() {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	fmt.Println(data)

	// compress
	var b bytes.Buffer
	w := lzma.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		fmt.Printf("lzma.Write() failed: %s", err.Error())
	}
	w.Close()
	fmt.Println(b.Bytes())
}
