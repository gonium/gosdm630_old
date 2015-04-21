package sdm630

import (
	"fmt"
)

type TextDumper struct {
	datastream ReadingChannel
	control    ControlChannel
}

func NewTextDumper(ds ReadingChannel, c ControlChannel) *TextDumper {
	return &TextDumper{datastream: ds, control: c}
}

func (td *TextDumper) ConsumeData() {
	for {
		readings := <-td.datastream
		fmt.Printf("%s\n", &readings)
	}
}
