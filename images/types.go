package images

type Image struct {
	Id         string
	Name       string
	Repository string
	Tag        string
	layers     []Layer
	ConfigFile 	  map[string]interface{}
}

type Layer struct {
	id      int64
	version string
	files   map[string][]byte
}


type ImageAnalyzer struct {
	JsonFiles    map[string][]byte
	Layers        map[string]map[string][]byte
}

func NewImageAnalyzer() *ImageAnalyzer{
	JsonFiles := make(map[string][]byte)
	Layers := make(map[string]map[string][]byte)
	return &ImageAnalyzer{JsonFiles, Layers}
}
