package filesystem

import (
	"archive/tar"
	"fmt"
	"io"
	"github.com/jeffthorne/beagle/images"
	"github.com/spf13/afero"
)

// ProcessLayerFileSystem processes a tar file associated with a container layer
// adding each file to a virtual filesystem. Virtual FS stored as part of Layer struct
// within Image struct.
func ProcessLayerFileSystem(layerDotTar *tar.Reader, layer *images.Layer) {

	for {
		header, err := layerDotTar.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error processing layer filesystem", err)
		}

		addToFileSystem(header, layer.FileSystem)
	}

}

func addToFileSystem(header *tar.Header, fs *afero.Fs) {

	if *fs == nil{
		fmt.Println()
	}
	err := afero.WriteFile(*fs, "/" + header.Name, make([]byte, header.Size), header.FileInfo().Mode())

	if err != nil {
		fmt.Println("Error writing file to virtual file system: ", err)
	}

}
