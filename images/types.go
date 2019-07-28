package images

type Image struct {
	Id         		string
	Name       		string
	Repository 		string
	Tag        		string
	ManifestFile 	[]byte
	ManifestJson	map[string]interface{}
	ConfigFile 		[]byte
	ConfigJson      map[string]interface{}
	Layers 			map[string]Layer
}

type Layer struct {
	Id      	string 	 //directory name
	Version		string
	Digest  	[32]byte   //sha256 on layer diff contents
	DigestString string
	Files   	map[string][]byte
	Author	    string
	CreatedBy   string
	Created     string
	Size        uint64
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
