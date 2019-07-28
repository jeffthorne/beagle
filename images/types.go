package images

type Image struct {
	Id         		string
	Name       		string
	Repository 		string
	Tag        		string
	ManifestFile 	[]byte
	ConfigFile 		[]byte
	Layers 			map[string]Layer
}

type Layer struct {
	Id      	string 	 //directory name
	Version		string
	Digest  	[32]byte   //sha256 on layer diff contents
	Files   	map[string][]byte
}

func NewLayer() *Layer{
	Files := make(map[string][]byte)
	return &Layer{Files:Files}
}

type ImageAnalyzer struct {
	Image		Image
	JsonFiles 	map[string][]byte
	Layers    	map[string]map[string][]byte
}

func NewImageAnalyzer() *ImageAnalyzer {
	Image := Image{}
	JsonFiles := make(map[string][]byte)
	Layers := make(map[string]map[string][]byte)
	return &ImageAnalyzer{Image, JsonFiles, Layers}
}
