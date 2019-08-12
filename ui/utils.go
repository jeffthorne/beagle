package ui

import (
	"bytes"
	"strconv"
	"strings"
	"fmt"
	"text/tabwriter"
	"github.com/jeffthorne/beagle/images"
)


const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)


func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "E "
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "P "
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "T "
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = " G "
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K "
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B "
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	result = fmt.Sprintf("%4s %s ", result, unit)
	return result
}


func ImageInfo(image *images.Image) string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(w, "\n   %s/%s:%s", image.Repository, image.Name, image.Tag)
	w.Flush()
	return b.String()
}

func LayerParagraph(layerNumber int, image *images.Image) string {

	if layerNumber < 1 {
		layerNumber = 1
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 1, '\t', tabwriter.AlignRight)
	imagesLayers := image.ManifestJson["Layers"].([]interface{})
	layer := imagesLayers[layerNumber-1]
	digest := strings.Split(layer.(string), "/")[0]
	l := image.Layers[digest]
	fmt.Fprintf(w, "\n[green]Digest\n%s\n\n[green]Command\n%s", l.DigestString, l.CreatedBy[11:])
	//b.WriteString("\n\t\tDigest: " + l.DigestString + "\n\n\tCommand:\n" + l.CreatedBy[11:])
	w.Flush()
	return b.String()

}


