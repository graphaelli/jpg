package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/graphaelli/jpg/structure"
)

// printable returns the printable chars in b up to length chars, kind of like strconv.QuoteToASCII
func printable(b []byte, length int) string {
	var clean []byte
	for i, c := range b {
		add := c
		if c < ' ' || c > 127 {
			add = '.'
		}
		clean = append(clean, add)
		if i >= length {
			break
		}
	}
	return string(clean)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s <filename.jpg>\n", os.Args[0])
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	j, err := structure.Load(bufio.NewReader(f))
	if err != nil && err != io.EOF {
		log.Fatalln("error", err)
	}
	fmt.Println(" address | marker     |  length | data")
	for _, m := range j {
		fmt.Printf("%8d | 0x%02x %-5s | %7d | %s\n",
			m.Address, m.Marker, structure.MarkerName(m.Marker), m.Length, printable(m.Data, 64))
	}
}
