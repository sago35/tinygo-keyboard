package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/itchio/lzma"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("this program needs vial.json's path")
	}
	// fmt.Println(os.Args[1])

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	r, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// compress
	var tBuf bytes.Buffer
	w := lzma.NewWriter(&tBuf)
	_, err = w.Write(r)
	if err != nil {
		fmt.Printf("lzma.Write() failed: %s", err.Error())
	}
	w.Close()

	var oBuf strings.Builder
	oBuf.WriteString("package main\n")
	oBuf.WriteString("\n")
	oBuf.WriteString("import keyboard \"github.com/sago35/tinygo-keyboard\"\n")
	oBuf.WriteString("\n")
	oBuf.WriteString("func loadKeyboardDef() {\n")
	oBuf.WriteString("\tkeyboard.KeyboardDef = []byte{\n")
	oBuf.WriteString("\t\t")
	for i, b := range tBuf.Bytes() {
		if i == (len(tBuf.Bytes()) - 1) {
			oBuf.WriteString(fmt.Sprintf("0x%02X,", b))
		} else {
			oBuf.WriteString(fmt.Sprintf("0x%02X, ", b))
		}
	}
	oBuf.WriteString("\n")
	oBuf.WriteString("\t}\n")
	oBuf.WriteString("}\n")

	outPath := filepath.Join(filepath.Dir(os.Args[1]), `def.go`)
	// fmt.Println(outPath)
	err = os.WriteFile(outPath, []byte(oBuf.String()), 0666)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(tBuf.Bytes())
}
