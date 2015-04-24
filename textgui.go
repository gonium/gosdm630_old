package sdm630

import (
	"bytes"
	ui "github.com/gizak/termui"
	"github.com/zfjagann/golang-ring"
	"text/template"
)

type TextGui struct {
	datastream ReadingChannel
	control    ControlChannel
	lastpower  *ring.Ring
}

const linedataTemplate = `
   {{.Voltage}} V
   {{.Current}} A
   {{.Power}} W
   {{.CosPhi}} cos phi
`

type linedata struct {
	Voltage float32
	Current float32
	Power   float32
	CosPhi  float32
}

func NewTextGui(ds ReadingChannel, c ControlChannel) *TextGui {
	r := &ring.Ring{}
	r.SetCapacity(600)
	return &TextGui{
		datastream: ds,
		control:    c,
		lastpower:  r,
	}
}

func formatStatsText(ld linedata) (retval string,
	err error) {
	t := template.Must(template.New("linedata").Parse(linedataTemplate))
	buf := new(bytes.Buffer)
	err = t.Execute(buf, ld)
	if err == nil {
		retval = buf.String()
	}
	return
}

func (td *TextGui) ConsumeData() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.UseTheme("helloworld")

	pl1 := ui.NewPar(" L1 ")
	pl1.Height = 8
	//p.Width = 50
	pl1.Border.Label = " L1 "

	pl2 := ui.NewPar(" L2 ")
	pl2.Height = 8
	//p.Width = 50
	pl2.Border.Label = " L2 "

	pl3 := ui.NewPar(" L3 ")
	pl3.Height = 8
	//p.Width = 50
	pl3.Border.Label = " L3 "

	lc0 := ui.NewLineChart()
	lc0.Border.Label = " Aggregated Load "
	//lc0.Data = td.lastl1voltage.Values()
	lc0.Data = []float64{0.1, 0.2, 0.3, 1024}
	lc0.Width = 80
	lc0.Height = 26
	lc0.X = 0
	lc0.Y = 0
	lc0.AxesColor = ui.ColorWhite
	lc0.LineColor = ui.ColorGreen | ui.AttrBold

	log := ui.NewPar(" Log ")
	log.Height = 6
	log.Border.Label = " Log "

	// Build the UI
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, pl1),
			ui.NewCol(4, 0, pl2),
			ui.NewCol(4, 0, pl3),
		),
		ui.NewRow(
			ui.NewCol(12, 0, lc0),
		),
		ui.NewRow(
			ui.NewCol(12, 0, log),
		),
	)

	draw := func(r Readings) {
		l1data := linedata{Voltage: r.L1Voltage,
			Current: r.L1Current,
			Power:   r.L1Power,
			CosPhi:  r.L1CosPhi,
		}
		pl1.Text, _ = formatStatsText(l1data)
		l2data := linedata{Voltage: r.L2Voltage,
			Current: r.L2Current,
			Power:   r.L2Power,
			CosPhi:  r.L2CosPhi,
		}
		pl2.Text, _ = formatStatsText(l2data)
		l3data := linedata{Voltage: r.L3Voltage,
			Current: r.L3Current,
			Power:   r.L3Power,
			CosPhi:  r.L3CosPhi,
		}
		pl3.Text, _ = formatStatsText(l3data)

		convValues := []float64{}
		for _, value := range td.lastpower.Values() {
			v := float64(value.(float32))
			convValues = append([]float64{v}, convValues...)
		}
		lc0.Data = convValues

		ui.Body.Align()
		ui.Render(ui.Body)
	}

	evt := ui.EventCh()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				td.control <- ControlClose
			}
		case readings := <-td.datastream:
			td.lastpower.Enqueue(readings.L1Power + readings.L2Power +
				readings.L3Power)
			draw(readings)
		}
	}
}
