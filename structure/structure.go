package structure

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

const (
	// Start of Frame markers, nondifferential Huffman-coding frames
	sof0Marker = 0xc0 // Start Of Frame (Baseline Sequential).
	/*
		sof1Marker = 0xc1 // Start Of Frame (Extended Sequential).
		sof2Marker = 0xc2 // Start Of Frame (Progressive).
		sof3Marker = 0xc3 // Start Of Frame (Lossless Sequential, Huffman coding).

		// Start of Frame markers, differential Huffman-coding frames
		sof5Marker = 0xc5 // Start Of Frame (differential sequential DCT, Huffman coding)
		sof6Marker = 0xc6 // Start Of Frame (differential progressive DCT, Huffman coding)
		sof7Marker = 0xc7 // Start Of Frame (differential lossless, Huffman coding)

		// Start of Frame markers, nondifferential arithmetic-coding frames
		sof9Marker  = 0xc9 // Start Of Frame (extended sequential DCT, arithmetic coding)
		sof10Marker = 0xca // Start Of Frame (progressive DCT, arithmetic coding)
		sof11Marker = 0xcb // Start Of Frame (lossless sequential, arithmetic coding)

		// Start of Frame markers, differential arithmetic-coding frames
		sof13Marker = 0xcd // Start Of Frame (differential sequential DCT, arithmetic coding)
		sof14Marker = 0xce // Start Of Frame (progressive DCT, arithmetic coding)
	*/
	sof15Marker = 0xcf // Start Of Frame (differential lossless, arithmetic coding)

	dhtMarker = 0xc4 // Define Huffman Table.
	soiMarker = 0xd8 // Start Of Image.
	eoiMarker = 0xd9 // End Of Image.
	sosMarker = 0xda // Start Of Scan.
	dqtMarker = 0xdb // Define Quantization Table.
	driMarker = 0xdd // Define Restart Interval.
	comMarker = 0xfe // COMment.
	// "APPlication specific" markers aren't part of the JPEG spec per se,
	// but in practice, their use is described at
	// http://www.sno.phy.queensu.ca/~phil/exiftool/TagNames/JPEG.html
	app0Marker  = 0xe0
	app15Marker = app0Marker | 0x0f
)

var markerName map[byte]string

type Mark struct {
	Address int
	Marker  byte
	Length  uint16
	Data    []byte
}

func init() {
	markerName = map[byte]string{
		soiMarker: "SOI",
		dhtMarker: "DHT",
		eoiMarker: "EOI",
		sosMarker: "SOS",
		dqtMarker: "DQT",
		driMarker: "DRI",
		comMarker: "COM",
	}
	var i byte
	for i = 0; i <= 15; i++ {
		markerName[app0Marker+i] = fmt.Sprintf("APP%d", i)
		if i != 4 { // % 4?
			markerName[sof0Marker+i] = fmt.Sprintf("SOF%d", i)
		}
	}
}

func MarkerName(marker byte) string {
	return markerName[marker]
}

// advanceToMarker moves j to the start of the next Marker, returning
// the number of bytes consumed and either the first byte of the next part or an error
func advanceToMarker(j *bufio.Reader) (int, byte, error) {
	var consumed int = 0
	var b byte
	var err error

	var stop byte = 0xff
	// Skips potential padding between markers
	for {
		b, err = j.ReadByte()
		if err != nil {
			return consumed, b, err
		}
		consumed++
		if b == stop {
			break
		}
	}

	// markers can start with any number of 0xff
	for {
		b, err = j.ReadByte()
		if err != nil {
			return consumed, b, err
		}
		consumed++
		if b != stop {
			break
		}
	}
	return consumed, b, nil
}

func hasLength(marker byte) bool {
	return (marker >= sof0Marker && marker <= sof15Marker) ||
		(marker >= app0Marker && marker <= app15Marker) ||
		marker == dhtMarker || marker == dqtMarker || marker == driMarker || marker == comMarker || marker == sosMarker
}

func Load(r *bufio.Reader) ([]Mark, error) {
	var marks []Mark
	var marker byte

	pos := 0
	for marker != eoiMarker && marker != sosMarker {
		if c, m, err := advanceToMarker(r); err != nil {
			return nil, err
		} else {
			pos += c // position advance
			marker = m
		}

		mark := Mark{
			Address: pos,
			Marker:  marker,
		}

		if hasLength(marker) {
			var length uint16
			if err := binary.Read(r, binary.BigEndian, &length); err != nil {
				return nil, err
			}
			pos += 2 // position advance
			mark.Length = length

			if marker >= app0Marker && marker <= app15Marker {
				b := make([]byte, length-2)
				if c, err := r.Read(b); err != nil {
					return nil, err
				} else {
					pos += c // position advance
				}
				mark.Data = b
			}
		}
		marks = append(marks, mark)
	}
	return marks, nil
}
