package images

import (
	"errors"
	"github.com/spf13/afero"
	"strings"
)

type Image struct {
	Id           string
	Name         string
	Repository   string
	Tag          string
	ManifestFile []byte
	ManifestJson map[string]interface{}
	ConfigFile   []byte
	ConfigJson   map[string]interface{}
	Layers       map[string]*Layer
}

func (i Image)InitialLayer() (string, error){

	layers := i.ManifestJson["Layers"].([]interface{})

	for k, v := range layers{
		if k == 0{
			str := strings.Split(v.(string), "/")[0]
			return str, nil
		}
	}

	return "", errors.New("could not find initial layer from manifest file")
}

type Layer struct {
	Id           string //directory name
	Version      string
	Digest       [32]byte //sha256 on layer diff contents
	DigestString string
	Files        map[string][]byte
	Author       string
	CreatedBy    string
	Created      string
	Size         uint64
	FileSystem   *afero.Fs
}

func NewLayer() *Layer {
	Files := make(map[string][]byte)
	FileSystem := afero.NewMemMapFs()
	return &Layer{Files: Files, FileSystem: &FileSystem}
}

type ImageAnalyzer struct {
	Image     Image
	JsonFiles map[string][]byte
	Layers    map[string]map[string][]byte
}

func NewImageAnalyzer() *ImageAnalyzer {
	Image := Image{}
	JsonFiles := make(map[string][]byte)
	Layers := make(map[string]map[string][]byte)
	return &ImageAnalyzer{Image, JsonFiles, Layers}
}
