package main

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jeffthorne/beagle/utils"
	dashboard "github.com/jeffthorne/beagle/ui"
	"log"
	"os"
	ui "github.com/gizak/termui/v3"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	fmt.Println("Processing image at:", filePath)
	image := utils.ProcessTar(filePath)
	for k, _ := range image.Layers{
		fmt.Printf("KEY:%s:\n", k)
	}
	listData := dashboard.Layers(image)

	if err := ui.Init(); err != nil{
		log.Fatalf("faile to initialize termui: %v", err)

	}

	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Layers"
	l.Rows = listData
	//l.SetRect(0, 0, 150, 50)
	l.TextStyle.Fg = ui.ColorWhite
	l.SelectedRowStyle.Fg = ui.ColorBlue
	l.SelectedRow = 1


//	draw := func(count int) {
	//	l.Rows = listData[:]
//		ui.Render(l)
	//}

	p := widgets.NewParagraph()
	p.Text = dashboard.LayerParagraph(l.SelectedRow, image)
	p.Title = "Layer Details"

	i := widgets.NewParagraph()
	i.Text = dashboard.ImageInfo(image)
	i.Title = "Image Info"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0,0, termWidth, termHeight)


	grid.Set(
		ui.NewRow(0.10,
		          ui.NewCol(0.6, i),
		),
		ui.NewRow(0.35,
			       ui.NewCol(0.6, l),
			       ),
		ui.NewRow(0.55,
			ui.NewCol(0.6,p),
			))

	ui.Render(grid)
//	tickerCount := 1
//	draw(tickerCount)

	previousKey := ""
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
			p.Text = dashboard.LayerParagraph(l.SelectedRow, image)
		case "k", "<Up>":
			l.ScrollUp()
			if l.SelectedRow < 1{
				l.SelectedRow = 1
			}
			p.Text = dashboard.LayerParagraph(l.SelectedRow, image)
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "<Enter>":
			p.Text = dashboard.LayerParagraph(l.SelectedRow, image)
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(grid)
	}







}