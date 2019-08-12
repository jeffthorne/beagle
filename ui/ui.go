package ui

import (
	"github.com/jeffthorne/beagle/dashboard"
	"github.com/jeffthorne/beagle/utils"
	"github.com/jeffthorne/beagle/images"
	"path/filepath"
	"bytes"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/afero"
	"strings"
	"text/tabwriter"
)


func SetupApp() *tview.Application{
	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(s tcell.Screen)bool{
		s.Clear()
		return false
	})

	return app
}

func SetupUIWidgets(image *images.Image, app *tview.Application) (*tview.TextView, *tview.TreeView, *tview.TextView, *tview.List, *tview.Flex){
	ld := tview.NewTextView()
	tv := tview.NewTreeView()
	id := tview.NewTextView()
	layersWidget := LayersWidget(image, ld, tv, app)

	LayerDetails(1, image, ld)
	ImageDetailsWidget(id,image)
	IntialFileSystem(image)
	inialLayer, _ := image.InitialLayer()
	LayerFileSystem(tv, image.Layers[inialLayer].FileSystem, layersWidget, app)

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(layersWidget, 0, 2, true).
			AddItem(ld, 0, 2, false).
			AddItem(id, 0, 1, false), 0, 2, false).
		AddItem(tv, 0, 2, false)
	flex.SetBackgroundColor(tcell.ColorDefault)

	return ld, tv, id, layersWidget, flex
}


func LayerFileSystem(tv *tview.TreeView, layerFileSystem *afero.Fs, ld *tview.List, app *tview.Application){

	rootDir := "/"
	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorRed)
	tv.SetRoot(root).SetCurrentNode(root)
	tv.SetTitle("Layer Details").SetTitleAlign(tview.AlignLeft)
	tv.SetBorder(true).SetBorderPadding(1,1,4,1)
	tv.SetBackgroundColor(tcell.ColorDefault)
	add(root, rootDir, layerFileSystem)
	tv.GetDrawFunc()
	tv.SetInputCapture(func(key *tcell.EventKey)*tcell.EventKey{
		if key.Name() == "Tab"{
			app.SetFocus(ld)
			tv.Blur()
			return nil

		}else {
			return key
		}
	})

	tv.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path, layerFileSystem)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

}


func add(target *tview.TreeNode, path string, fs *afero.Fs)  {

	files, err := afero.ReadDir(*fs, path)

	if err != nil {
		panic(err)
	}
	for _, file := range files {
		//fileName :=  fmt.Sprintf("%s     %s", file.Mode().String(), file.Name())
		var b bytes.Buffer
		w := tabwriter.NewWriter(&b, 13, 1, 2, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "%s\t   %5d %s", file.Mode().String(), file.Size(), file.Name())
		w.Flush()


		node := tview.NewTreeNode(b.String()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(file.IsDir())
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}

}




func IntialFileSystem(image *images.Image){
	manifestLayers := image.Layers

	for k, _ := range manifestLayers{
		fmt.Println("LAYER ID: ", k)
	}
}

func LayersWidget(image *images.Image, ld *tview.TextView, tv *tview.TreeView, app *tview.Application) *tview.List{
	box := tview.NewBox()
	box.SetBorder(true)

	box.SetTitle("Layers")
	box.SetTitleAlign(tview.AlignLeft)
	list := tview.NewList()

	layerText := Layers(image)

	for _, text := range layerText {
		list.AddItem(text, "Some explanatory text", 0, nil)
	}

	list.ShowSecondaryText(false)
	list.SetTitle("Layers").SetTitleAlign(tview.AlignLeft)
	list.SetBorder(true).SetBorderPadding(1,0,0,0)
	list.SetShortcutColor(tcell.ColorDefault)
	list.SetSelectedTextColor(tcell.ColorGreen)
	list.SetSelectedBackgroundColor(tcell.ColorDefault)
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetCurrentItem(1)
	list.SetInputCapture(func(key *tcell.EventKey)*tcell.EventKey{
		//fmt.Println("KEY RECEIVED: ", key.Name())
		if key.Name() == "Tab"{
			//fmt.Println("Changing focus")
			//fmt.Println(app.GetFocus()) /
			app.SetFocus(tv)
			list.Blur()
			return nil

		}else {
			return key
		}
	})

	list.SetChangedFunc(func(index int, tableName string, t string, s rune){
		if index == 0{
			index = 1
			list.SetCurrentItem(index)
		}


		LayerDetails(index, image, ld)
		imagesLayers := image.ManifestJson["Layers"].([]interface{})
		layer := imagesLayers[index-1]
		digest := strings.Split(layer.(string), "/")[0]
		l := image.Layers[digest]
		LayerFileSystem(tv, l.FileSystem, list, app)
	})

	return list
}


func Layers(image *images.Image) []string {

	imagesLayers1 := image.ManifestJson["Layers"].([]interface{})
	var layers []string

	count := 1
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 7, 1, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "\t%v\t     %v\t", "Size", "Command")
	w.Flush()
	layers = append(layers, b.String())
	b.Reset()

	for _, il := range imagesLayers1 {
		digest := strings.Split(il.(string), "/")[0]
		fmt.Printf("DIGEST:%s:\n", digest)
		if _, ok := image.Layers[digest]; ok {
			fmt.Printf("%d - %sin List:", count, digest)
		}

		l := image.Layers[digest]
		bs := dashboard.ByteSize(l.Size)

		fmt.Fprintf(w, "\t%s\t   %s\t", bs, utils.StringMaxSize(l.CreatedBy[11:], 45))
		w.Flush()
		layers = append(layers, b.String())
		b.Reset()
		count++
	}

	return layers
}

func ImageDetailsWidget(id *tview.TextView, image *images.Image){
	id.SetTitle("Image Details")
	id.SetBackgroundColor(tcell.ColorDefault)
	id.SetBorder(true)
	id.SetTitleAlign(tview.AlignLeft)
	id.SetText(dashboard.ImageInfo(image))

}



func LayerDetails(selectedRow int, image *images.Image, ld *tview.TextView){
	ld.SetDynamicColors(true)
	ld.SetText(dashboard.LayerParagraph(selectedRow, image))
	ld.SetBorder(true)
	ld.SetTitle("Layer Details")
	ld.SetTitleAlign(tview.AlignLeft)
	ld.SetBorderPadding(1,1,1,1)
	ld.SetScrollable(true)
	ld.SetBackgroundColor(tcell.ColorDefault)
	ld.SetWrap(true)
}

