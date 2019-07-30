package main

import (
	"fmt"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jeffthorne/beagle/dashboard"
	"github.com/jeffthorne/beagle/utils"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	fmt.Println("Processing image at:", filePath)
	image := utils.ProcessTar(filePath)

	if err := ui.Init(); err != nil {
		log.Fatalf("faile to initialize termui: %v", err)
	}
	defer ui.Close()

	termHeight := 0

	l := dashboard.UILayersWidget(image, widgets.NewList())
	p := dashboard.UILayerDetailsWidget(l.SelectedRow, image, widgets.NewParagraph())
	i := dashboard.UIImageDetailsWidget(image, termHeight)
	bar := dashboard.NewStatusBar(image)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight-1)

	grid.Set(
		ui.NewRow(0.10,
			ui.NewCol(0.6, i),
		),
		ui.NewRow(0.35,
			ui.NewCol(0.6, l),
		),
		ui.NewRow(0.55,
			ui.NewCol(0.6, p),
		))

	ui.Render(grid)
	bar.SetRect(0, termHeight-1, termWidth, termHeight)
	ui.Render(bar)

	previousKey := ""
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
			p.Text = dashboard.UILayerDetailsWidget(l.SelectedRow, image, p).Text
		case "k", "<Up>":
			l.ScrollUp()
			if l.SelectedRow < 1 {
				l.SelectedRow = 1
			}
			p.Text = dashboard.UILayerDetailsWidget(l.SelectedRow, image, p).Text
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "<Enter>":
			p.Text = dashboard.UILayerDetailsWidget(l.SelectedRow, image, p).Text
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()

		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			//grid.SetRect(0, 0, payload.Width, payload.Height)
			termWidth, termHeight := payload.Width, payload.Height
			grid.SetRect(0, 0, payload.Width, payload.Height-1)
			//ui.Clear()
			ui.Render(grid)
			bar.SetRect(0, termHeight-1, termWidth, termHeight)
			ui.Render(bar)
		}
		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(grid)
	}

}
