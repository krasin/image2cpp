package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Interval struct {
	DataOff int
	Off     int
	Len     int
}

type Image struct {
	Data    []byte
	Sectors map[int][]Interval
}

func NewImage() *Image {
	return &Image{
		Sectors: make(map[int][]Interval),
	}
}

func (im *Image) Add(sector int, data []byte) {
	var intervals []Interval
	ind := -1
	off := 0
	zcount := 0
	for i, b := range data {
		if b != 0 {
			if zcount > 10 && ind >= 0 {
				dataOff := len(im.Data)
				im.Data = append(im.Data, data[off:ind+1]...)
				intervals = append(intervals, Interval{dataOff, off, ind + 1 - off})
				off = i
			}
			if zcount > 0 && ind == -1 {
				off = i
			}
			ind = i
			zcount = 0
		} else {
			zcount++
		}

	}
	if ind == -1 {
		return
	}

	dataOff := len(im.Data)
	im.Data = append(im.Data, data[off:ind+1]...)
	intervals = append(intervals, Interval{dataOff, off, ind + 1 - off})
	im.Sectors[sector] = intervals
}

func main() {
	im := NewImage()
	buf := make([]byte, 512)
	for i := 0; ; i++ {
		_, err := io.ReadFull(os.Stdin, buf)
		if err == io.EOF {
			break
		}
		if err == io.ErrUnexpectedEOF {
			log.Fatalf("Unexpected EOF from input stream")
		}
		if err != nil {
			log.Fatalf("ReadFull: %v", err)
		}
		im.Add(i, buf)
	}

	fmt.Printf("char image_data[] = {\n    ")
	x := 0
	for i, b := range im.Data {
		if i > 0 {
			fmt.Printf(", ")
			x += 2
			if x > 70 {
				x = 0
				fmt.Printf("\n    ")
			}
		}
		switch {
		/*				case b >= 'a' && b <= 'z' ||
							b >= 'A' && b <= 'Z' ||
							b >= '0' && b <= '9' ||
							b == ' ' || b == '.' || b == ',':
							x++
							fmt.Printf("%c", b)
						case b < 10:
							x += 2
							fmt.Printf("\\%d", b)
						case b == 10:
							x += 2
							fmt.Printf("\\n")
						case b == 13:
							x += 2
							fmt.Printf("\\r")*/
		case b >= 'a' && b <= 'z' ||
			b >= 'A' && b <= 'Z' ||
			b >= '0' && b <= '9' ||
			b == ' ' || b == '.' || b == ',':
			x += 3
			fmt.Printf("'%c'", b)
		case b < 10:
			x--
			fallthrough
		case b < 100:
			x--
			fallthrough
		default:
			x += 3
			fmt.Printf("%d", b)
		}
	}
	fmt.Printf("};\n\n")

	fmt.Printf("void copy_sector_data(char* buffer, int block_number) {\n")
	fmt.Printf("  switch (block_number) {\n")
	for block_number, intervals := range im.Sectors {
		fmt.Printf("  case %d:\n", block_number)
		for _, interval := range intervals {
			fmt.Printf("    memcpy(buffer+%d, image_data+%d, %d);\n", interval.Off, interval.DataOff, interval.Len)
		}
		fmt.Printf("    break;\n")
	}
	fmt.Printf("  }\n")
	fmt.Printf("}\n\n")
}
