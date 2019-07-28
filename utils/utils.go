package utils

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
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

func ProcessTar(filepath string) images.Image{

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
			continue
		case tar.TypeSymlink, tar.TypeLink, tar.TypeReg:
			layerId := ""
			f, err := ioutil.ReadAll(tr)
			if err != nil{
				fmt.Println(err)
			}

			if strings.Contains(filename, "json") {

				if filename == "manifest.json" {
					ParseManifest(f, &imageAnalyer.Image)
					layerId = imageAnalyer.Image.Id
					imageAnalyer.JsonFiles[filename] = f
				} else if strings.Contains(filename, ".") {
					ParseConfig(f, &imageAnalyer.Image)
					layerId = strings.Split(filename, ".")[0]
					imageAnalyer.JsonFiles[layerId+".json"] = f
				} else {
					layerId = getLayerId(f)
					imageAnalyer.JsonFiles[layerId+".json"] = f
				}


			}else{

				if filename != "repositories" {
					layerId = strings.Split(filename, "/")[0]
					filename = strings.Split(filename, "/")[1]
				}
			}

			if layerId != "" {
				if imageAnalyer.Layers[layerId] == nil {
					imageAnalyer.Layers[layerId] = make(map[string][]byte)
				}

				if strings.Contains(filename, "/json"){
					filename = strings.Split(filename, "/")[1]
				}
				imageAnalyer.Layers[layerId][filename] = f


			}



		}

		if err != nil {
			fmt.Println(err)
		}
	}

	return makeImageStruct(imageAnalyer)

}

func getLayerId(f []byte) string {

	var result map[string]interface{}
	json.Unmarshal(f, &result)
	return result["id"].(string)
}

func ParseConfig(f []byte, image *images.Image) {

	var result map[string]interface{}
	json.Unmarshal(f, &result)
	image.ConfigFile = f
}

func ParseManifest(f []byte, image *images.Image) {

	var result []map[string]interface{}

	json.Unmarshal(f, &result)
	image.Id = strings.Split(result[0]["Config"].(string), ".")[0]
	repoTags := result[0]["RepoTags"].([]interface{})[0].(string)
	if strings.Contains(repoTags, "/") {
		image.Repository = strings.Split(repoTags, "/")[0]
		image.Name = strings.Split(strings.Split(repoTags, "/")[1], ":")[0]
	}else{
		image.Repository = ""
		image.Name = strings.Split(strings.Split(repoTags, ":")[0], ":")[0]
	}

	image.Tag = strings.Split(repoTags, ":")[1]

}


func makeImageStruct(ia *images.ImageAnalyzer) images.Image{
	image := ia.Image

	for k, v := range ia.Layers{

		if _, ok := v["manifest.json"]; ok{
			image.Id = k
			image.ManifestFile = v["manifest.json"]
		}else if _, ok := v["layer.tar"]; ok{
			layer := images.Layer{}
			layer.Files = v
			layer.Digest = sha256.Sum256(v["layer.tar"])
			digestString := hex.EncodeToString(layer.Digest[:])
			if image.Layers == nil{
				image.Layers = make(map[string]images.Layer)
			}
			image.Layers[digestString] = layer
		}


	}

	return image

}
