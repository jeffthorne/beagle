package main

import (
	"github.com/jeffthorne/beagle/utils"
	"github.com/jeffthorne/beagle/ui"
	"os"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	app := ui.SetupApp()
	image := utils.ProcessTar(filePath)

	_, _, _, layersWidget, flex := ui.SetupUIWidgets(image, app)  //TODO: need to refactor UI code


	if err := app.SetRoot(flex, true).SetFocus(layersWidget).Run(); err != nil {
		panic(err)
	}
}


