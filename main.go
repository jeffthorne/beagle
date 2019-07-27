package main

import (
	"fmt"
	"os"

	"github.com/jeffthorne/beagle/utils"
)

// main takes on command argument. Full path to container .tar
func main() {

	filePath := os.Args[1]
	fmt.Println(filePath)

	utils.ProcessTar(filePath)
}
