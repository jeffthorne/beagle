package utils

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jeffthorne/beagle/images"
)

var (
	readline    = 6
	gzipHeader  = []byte{0x1f, 0x8b}
	bzip2Header = []byte{0x42, 0x5a, 0x68}
	xzHeader    = []byte{0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00}
)

func ProcessTar(filepath string) {

	image := images.Image{}
	tarFile, err := os.Open(filepath) //Open the tar file
	defer tarFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	tr := tar.NewReader(tarFile)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break //end of archive
		}

		if err != nil {
			fmt.Println(err)
		}

		filename := header.Name
		filename = strings.TrimPrefix(filename, "./")

		switch header.Typeflag {
		case tar.TypeDir:
			fmt.Printf("%s is a directory\n", filename)
			continue
		case tar.TypeSymlink, tar.TypeLink, tar.TypeReg:
			fmt.Printf("%s is a file\n", filename)
			//d, _ := ioutil.ReadAll(tr)
			//data[filename] = d
			if filename == "manifest.json" {
				ParseManifest(tr, &image)
			}

		default:

		}

		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("Image Id is %s\n", image.Id)
	fmt.Println(image)

}

func ParseManifest(manifestFile *tar.Reader, image *images.Image) {

	var result []map[string]interface{}
	f, err := ioutil.ReadAll(manifestFile)

	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(f, &result)
	image.Id = strings.Split(result[0]["Config"].(string), ".")[0]
	repoTags := result[0]["RepoTags"].([]interface{})[0].(string)
	fmt.Println(repoTags)
	image.Repository = strings.Split(repoTags, "/")[0]
	image.Name = strings.Split(strings.Split(repoTags, "/")[1], ":")[0]
	image.Tag = strings.Split(repoTags, ":")[1]

}
