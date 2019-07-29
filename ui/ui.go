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
			str += "-"
		}
	}


	return str
}

func Layers(image *images.Image) []string{

	imagesLayers1 := image.ManifestJson["Layers"].([]interface{})
	var layers []string


	count := 1
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 7,4,1,' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "   \t%v\t\t%v\t\t\t\t\t\t\t   %s", "Size", "Command", "Digest")
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


		fmt.Fprintf(w, "[%d]\t%s\t\t%s\t\t%s", count, humanize.Bytes(l.Size), StringMaxSize(l.CreatedBy[11:], 45), l.DigestString)
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