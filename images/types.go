package images

type Image struct {
	Id         string
	Name       string
	Repository string
	Tag        string
	layers     []Layer
}

type Layer struct {
	id      int64
	version string
	files   map[string][]byte
}
