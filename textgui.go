package sdm630

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/zfjagann/golang-ring"
	"log"
)

type TextGui struct {
	datastream    ReadingChannel
	control       ControlChannel
	lastl1voltage *ring.Ring
}

func NewTextGui(ds ReadingChannel, c ControlChannel) *TextGui {
	r := &ring.Ring{}
	r.SetCapacity(10)
	return &TextGui{
		datastream:    ds,
		control:       c,
		lastl1voltage: r,
	}
}

func (td *TextGui) ConsumeData() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = "Text Box"
	p.Border.FgColor = ui.ColorCyan

	lc0 := ui.NewLineChart()
	lc0.Border.Label = "braille-mode Line Chart"
	//lc0.Data = td.lastl1voltage.Values()
	lc0.Data = []float64{0.1, 0.2, 0.3}
	lc0.Width = 50
	lc0.Height = 12
	lc0.X = 0
	lc0.Y = 0
	lc0.AxesColor = ui.ColorWhite
	lc0.LineColor = ui.ColorGreen | ui.AttrBold

	draw := func() {

		//lc0.Data = float64(td.lastl1voltage.Values())
		ui.Render(p, lc0)
	}

	evt := ui.EventCh()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				// TODO: Properly terminate the client
				log.Fatal("exiting.")
			}
		case readings := <-td.datastream:
			td.lastl1voltage.Enqueue(readings.L1Voltage)
			fmt.Printf("R: %s\n", &readings)
			draw()
		}
	}
}
