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
	imageAnalyer := images.NewImageAnalyzer()

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

			if strings.Contains(filename, "json"){
				layerId := ""

				if f, err := ioutil.ReadAll(tr); err == nil {
					if filename == "manifest.json"{
						ParseManifest(f, &image)
						layerId = image.Id
						imageAnalyer.JsonFiles[filename] = f
					}else if strings.Contains(filename, ".") {
						ParseConfig(f, &image)
						layerId = strings.Split(filename, ".")[0]
						imageAnalyer.JsonFiles[layerId + ".json"] = f
					}else{
						layerId = getLayerId(f)
						imageAnalyer.JsonFiles[layerId + ".json"] = f
					}

					if imageAnalyer.Layers[layerId] == nil{
						imageAnalyer.Layers[layerId] = make(map[string][]byte)
					}
					imageAnalyer.Layers[layerId][filename] = f
				}
			}



		default:

		}

		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("Image Id is %s\n", image.Id)


}

func getLayerId(f []byte) string{

	var result map[string]interface{}
	json.Unmarshal(f, &result)
	return result["id"].(string)
}

func ParseConfig(f []byte, image *images.Image){

	var result map[string]interface{}
	json.Unmarshal(f, &result)
	image.ConfigFile = result
}

func ParseManifest(f []byte, image *images.Image) {

	var result []map[string]interface{}

	json.Unmarshal(f, &result)
	image.Id = strings.Split(result[0]["Config"].(string), ".")[0]
	repoTags := result[0]["RepoTags"].([]interface{})[0].(string)
	image.Repository = strings.Split(repoTags, "/")[0]
	image.Name = strings.Split(strings.Split(repoTags, "/")[1], ":")[0]
	image.Tag = strings.Split(repoTags, ":")[1]

}
