package main

import (
	"fmt"
	"github.com/jeffthorne/beagle/utils"
	"os"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	fmt.Println("Processing image at:", filePath)
	image := utils.ProcessTar(filePath)

	for k, v := range image.Layers{
		fmt.Println(k)

		for k, _ := range v.Files {
			fmt.Println(k)
		}
		fmt.Println("\n\n")
	}


}
