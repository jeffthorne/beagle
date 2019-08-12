package ui

import (
	"image"

	ui2 "github.com/gizak/termui/v3"
	"github.com/jeffthorne/beagle/images"
)

type StatusBar struct {
	ui2.Block
	Image *images.Image
}

func NewStatusBar(image *images.Image) *StatusBar {
	self := &StatusBar{*ui2.NewBlock(), image}
	self.Border = false
	return self
}

func (self *StatusBar) Draw(buf *ui2.Buffer) {
	self.Block.Draw(buf)

	buf.SetString(
		//"Image: "+self.Image.Repository+"/"+self.Image.Name+":"+self.Image.Tag,
		"Image Format: Docker",
		ui2.NewStyle(ui2.ColorWhite),
		image.Pt(self.Inner.Min.X, self.Inner.Min.Y+(self.Block.Inner.Dy()/2)),
	)

}
