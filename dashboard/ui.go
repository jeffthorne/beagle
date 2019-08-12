package dashboard

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	//"github.com/dustin/go-humanize"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jeffthorne/beagle/images"
	"github.com/jeffthorne/beagle/utils"

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

var styleBold = termui.Style{termui.ColorWhite, termui.ColorClear, termui.ModifierBold}
type nodeValue string
func (nv nodeValue) String() string {
	return string(nv)
}

func UITreeWidget(image *images.Image) *widgets.Tree{
	t := widgets.NewTree()
	t.Title = "[Layer Contents]"
	t.TitleStyle = styleBold
	t.WrapText = false
	t.TextStyle.Fg = termui.ColorWhite
	t.BorderTop, t.BorderBottom, t.BorderLeft, t.BorderRight = true, false, true, false
	nodes := []*widgets.TreeNode{
		{
			Value: nodeValue("Key 1"),
			Nodes: []*widgets.TreeNode{
				{
					Value: nodeValue("Key 1.1"),
					Nodes: []*widgets.TreeNode{
						{
							Value: nodeValue("Key 1.1.1"),
							Nodes: nil,
						},
						{
							Value: nodeValue("Key 1.1.2"),
							Nodes: nil,
						},
					},
				},
				{
					Value: nodeValue("Key 1.2"),
					Nodes: nil,
				},
			},
		},
		{
			Value: nodeValue("Key 2"),
			Nodes: []*widgets.TreeNode{
				{
					Value: nodeValue("Key 2.1"),
					Nodes: nil,
				},
				{
					Value: nodeValue("Key 2.2"),
					Nodes: nil,
				},
				{
					Value: nodeValue("Key 2.3"),
					Nodes: nil,
				},
			},
		},
		{
			Value: nodeValue("Key 3"),
			Nodes: nil,
		},
	}
	t.SetNodes(nodes)
	return t
}

func UILayersWidget(image *images.Image, l *widgets.List) *widgets.List {
	l.Title = "[Layers]"
	l.TitleStyle = styleBold
	l.Rows = Layers(image)
	//l.SetRect(0, 0, 150, 50)
	l.TextStyle.Fg = termui.ColorWhite
	l.SelectedRowStyle.Fg = termui.ColorGreen
	l.SelectedRow = 1
	l.BorderTop, l.BorderBottom, l.BorderLeft, l.BorderRight = true, false, false, false

	return l
}

func UILayerDetailsWidget(selectedRow int, image *images.Image, p *widgets.Paragraph) *widgets.Paragraph {
	p.Text = LayerParagraph(selectedRow, image)
	p.Title = "[Layer Details]"
	p.TitleStyle = styleBold
	p.BorderTop, p.BorderBottom, p.BorderLeft, p.BorderRight = true, false, false, false
	return p
}

func UIImageDetailsWidget(image *images.Image, termHeight int) *widgets.Paragraph {
	i := widgets.NewParagraph()
	i.Text = ImageInfo(image)
	i.Title = "[Image Details]"
	i.TitleStyle = styleBold
	i.BorderTop, i.BorderBottom, i.BorderLeft, i.BorderRight = true, false, false, false
	return i
}

func Layers(image *images.Image) []string {

	imagesLayers1 := image.ManifestJson["Layers"].([]interface{})
	var layers []string

	count := 1
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 7, 1, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "\t%v\t       %v\t", "[Size](mod:bold)", "[Command](mod:bold)")
	w.Flush()
	layers = append(layers, b.String())
	b.Reset()

	for _, il := range imagesLayers1 {
		digest := strings.Split(il.(string), "/")[0]

		if _, ok := image.Layers[digest]; ok {
			fmt.Printf("%d - %sin List:", count, digest)
		}

		l := image.Layers[digest]

		//bc := utils.ByteCountIEC(int64(l.Size))
		//fmt.Println("BYTE COUNT****************************:", bc)

		bs := ByteSize(l.Size)

		fmt.Fprintf(w, "\t%s\t   %s\t", bs, utils.StringMaxSize(l.CreatedBy[11:], 45))
		w.Flush()
		layers = append(layers, b.String())
		b.Reset()
		count++
	}

	for _, v := range layers {
		fmt.Println("FINAL ITEM:", v)
	}

	return layers
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
