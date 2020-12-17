// tiledback
package tiledback

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

type tile struct {
	image.Image
	unitw, unith int
}

func newTile(img image.Image) *tile {
	r := new(tile)
	r.Image = img
	r.unitw = img.Bounds().Dx()
	r.unith = img.Bounds().Dy()
	return r
}

func (t *tile) genRaster(x, y, w, h int) color.Color {
	ux := x % t.unitw
	uy := y % t.unith
	return t.Image.At(ux, uy)
}

type TileBackground struct {
	widget.BaseWidget
	unitImage *tile
}

func NewTileBackground(unitImg image.Image) *TileBackground {
	tile := &TileBackground{
		unitImage: newTile(unitImg),
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
	genImg           *canvas.Raster
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
		tbr.genImg = canvas.NewRasterWithPixels(tbr.tb.unitImage.genRaster)
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
