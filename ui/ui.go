package ui

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/jeffthorne/beagle/images"
	"github.com/rivo/tview"
	"strings"
	"text/tabwriter"
)

func Layers(image *images.Image) *tview.List{

	layers := tview.NewList().ShowSecondaryText(false)
	layers.SetBorder(true).SetTitle("Layers")

	imagesLayers1 := image.ManifestJson["Layers"].([]interface{})
	//fmt.Println("******************",strings.Split(imagesLayers1[0].(string), "/")[0])
	//imagesLayers := image.Layers



	for _, il := range imagesLayers1 {
		digest := strings.Split(il.(string), "/")[0]
		l := image.Layers[digest]
		var b bytes.Buffer
		w := tabwriter.NewWriter(&b, 7,1,1,' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "%v\t\t%v\t", humanize.Bytes(l.Size), l.CreatedBy[:20])
		w.Flush()
		layers.AddItem(b.String(), "Some explanatory text",0, nil)
	}

	/*layers.AddItem("List item 1", "Some explanatory text", 'a', nil).
		AddItem("List item 2", "Some explanatory text", 'b', nil).
		AddItem("List item 3", "Some explanatory text", 'c', nil).
		AddItem("List item 4", "Some explanatory text", 'd', nil)
	*
	 */

	return layers
}