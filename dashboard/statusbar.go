package dashboard

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/jeffthorne/beagle/images"
)

type StatusBar struct {
	ui.Block
	Image *images.Image
}

func NewStatusBar(image *images.Image) *StatusBar {
	self := &StatusBar{*ui.NewBlock(), image}
	self.Border = false
	return self
}

func (self *StatusBar) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	buf.SetString(
		"Image: "+self.Image.Repository+"/"+self.Image.Name+":"+self.Image.Tag,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(self.Inner.Min.X, self.Inner.Min.Y+(self.Block.Inner.Dy()/2)),
	)

}
