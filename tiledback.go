// tiledback
package tiledback

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func TileImage(img image.Image, width, height int) image.Image {
	unitw := img.Bounds().Dx()
	unith := img.Bounds().Dy()
	if img.Bounds().Dx() >= width && img.Bounds().Dy() >= height {
		return img
	}
	curPoint := image.Pt(0, 0)
	rect := image.Rect(0, 0, width, height)
	rimg := image.NewRGBA(rect)
	var neww, newh int
	for hleft := height; hleft > 0; hleft = hleft - img.Bounds().Dy() {
		for wleft := width; wleft > 0; wleft = wleft - img.Bounds().Dx() {
			neww = unitw
			if width-curPoint.X < unitw {
				neww = width - curPoint.X
			}
			newh = unith
			if height-curPoint.Y < unith {
				newh = height - curPoint.Y
			}
			draw.Draw(rimg, image.Rect(curPoint.X, curPoint.Y, curPoint.X+neww, curPoint.Y+newh), img, image.Pt(0, 0), draw.Over)
			curPoint = image.Pt(curPoint.X+neww, curPoint.Y)
		}
		curPoint = image.Pt(0, curPoint.Y+newh)
	}
	return rimg
}

type TileBackground struct {
	widget.BaseWidget
	unitImage image.Image
}

func NewTileBackground(unitImg image.Image) *TileBackground {
	tile := &TileBackground{
		unitImage: unitImg,
	}
	tile.ExtendBaseWidget(tile)
	return tile
}

func NewTileBackgroundFromFile(filepath string) (*TileBackground, error) {
	freader, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(freader)
	if err != nil {
		return nil, err
	}

	return NewTileBackground(img), nil
}

func (tb *TileBackground) CreateRenderer() fyne.WidgetRenderer {
	tb.ExtendBaseWidget(tb)
	return newTileBackgroundRender(tb)
}

type tileBackgroundRender struct {
	tb               *TileBackground
	overallContainer *fyne.Container
	genImg           *canvas.Image
	mux              *sync.RWMutex
	curSize          fyne.Size
}

func newTileBackgroundRender(tb *TileBackground) *tileBackgroundRender {
	r := &tileBackgroundRender{
		tb:               tb,
		mux:              new(sync.RWMutex),
		overallContainer: fyne.NewContainerWithoutLayout(canvas.NewImageFromImage(tb.unitImage)),
	}
	return r
}
func (tbr *tileBackgroundRender) BackgroundColor() color.Color {
	return color.Transparent
}
func (tbr *tileBackgroundRender) Destroy() {

}

func (tbr *tileBackgroundRender) Layout(layoutsize fyne.Size) {
	tbr.mux.Lock()
	defer tbr.mux.Unlock()
	if len(tbr.overallContainer.Objects) == 0 || !layoutsize.Subtract(tbr.curSize).IsZero() {
		brimg := TileImage(tbr.tb.unitImage, layoutsize.Width, layoutsize.Height)
		tbr.genImg = canvas.NewImageFromImage(brimg)
		tbr.genImg.FillMode = canvas.ImageFillContain
		tbr.genImg.Resize(layoutsize)
		tbr.overallContainer = fyne.NewContainerWithoutLayout()
		tbr.overallContainer.Add(tbr.genImg)
		tbr.curSize = layoutsize
	}
}

func (tbr *tileBackgroundRender) MinSize() fyne.Size {
	return fyne.NewSize(tbr.tb.unitImage.Bounds().Dx(), tbr.tb.unitImage.Bounds().Dy())
}

func (tbr *tileBackgroundRender) Objects() []fyne.CanvasObject {
	tbr.mux.RLock()
	defer tbr.mux.RUnlock()
	return []fyne.CanvasObject{tbr.overallContainer}
}
func (tbr *tileBackgroundRender) Refresh() {
	canvas.Refresh(tbr.tb)
}
