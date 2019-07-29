package ui

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jeffthorne/beagle/images"
	"strings"
	"text/tabwriter"
     "unicode/utf8"
)

func StringMaxSize(str string, size int) string{

	str = strings.ReplaceAll(str, "\t", "  ")
	if(utf8.RuneCountInString(str) > size ){
		str = str[:size]
	}else{
		diff := size - utf8.RuneCountInString(str)

		for i := 0; i < diff; i++{
			str += " "
		}
	}


	return str
}

func Layers(image *images.Image) []string{

	imagesLayers1 := image.ManifestJson["Layers"].([]interface{})
	var layers []string


	count := 1
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 7,1,1,' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "   \t%v\t\t%v\t", "Size", "Command")
	w.Flush()
	layers = append(layers, b.String())
	b.Reset()

	for _, il := range imagesLayers1 {
		digest := strings.Split(il.(string), "/")[0]
		fmt.Printf("DIGEST:%s:\n", digest)
		if _, ok := image.Layers[digest]; ok{
			fmt.Printf("%d - %sin List:", count, digest)
		}

		l := image.Layers[digest]


		fmt.Fprintf(w, "[%d]\t%s\t\t%s\t", count, humanize.Bytes(l.Size), StringMaxSize(l.CreatedBy[11:], 45))
		w.Flush()
		layers = append(layers, b.String())
		b.Reset()
		count++
	}



	for _, v := range layers{
		fmt.Println("FINAL ITEM:", v)
	}


	return layers
}

func LayerParagraph(layerNumber int, image *images.Image) string{

	var b bytes.Buffer
	imagesLayers := image.ManifestJson["Layers"].([]interface{})
	layer := imagesLayers[layerNumber - 1]
	digest := strings.Split(layer.(string), "/")[0]
	l := image.Layers[digest]

	b.WriteString("\nDigest: " + l.DigestString + "\n\nCommand:\n" + l.CreatedBy[11:])
	return b.String()

}