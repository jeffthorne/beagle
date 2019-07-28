package main

import (
	"fmt"
	"github.com/jeffthorne/beagle/ui"
	"github.com/jeffthorne/beagle/utils"
	"github.com/rivo/tview"
	"os"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	fmt.Println("Processing image at:", filePath)
	image := utils.ProcessTar(filePath)


	app := tview.NewApplication()

	flex := tview.NewFlex().
		AddItem(ui.Layers(image), 0, 1, true).
		AddItem(tview.NewList(), 0, 1, false)

	app.SetRoot(flex, true).Run()




}