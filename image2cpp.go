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

	fmt.Printf("char image_data[] = {")
	for i, b := range im.Data {
		if i > 0 {
			fmt.Printf(", ")
			if i%30 == 0 {
				fmt.Printf("\n")
			}
		}
		fmt.Printf("%x", b)
	}
	fmt.Printf("};\n")
	fmt.Printf("Image.Sectors: %v\n", im.Sectors)
}
