// tiledback
package tiledback

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"

	// "fyne.io/fyne/theme"

	// "fyne.io/fyne/layout"
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
	log.Printf("creating new redener")
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
	log.Printf("layout for %v", layoutsize)
	if len(tbr.overallContainer.Objects) == 0 || !layoutsize.Subtract(tbr.curSize).IsZero() {
		log.Printf("working for %v", layoutsize)
		brimg := TileImage(tbr.tb.unitImage, layoutsize.Width, layoutsize.Height)
		log.Printf("created image with size %v", brimg.Bounds().Size())
		tbr.genImg = canvas.NewImageFromImage(brimg)
		tbr.genImg.Resize(layoutsize)
		tbr.overallContainer = fyne.NewContainerWithoutLayout()
		tbr.overallContainer.Add(tbr.genImg)
		tbr.curSize = layoutsize
		// tbr.overallContainer.Resize(layoutsize)

	}

}

func (tbr *tileBackgroundRender) MinSize() fyne.Size {
	return fyne.NewSize(tbr.tb.unitImage.Bounds().Dx(), tbr.tb.unitImage.Bounds().Dy())
}

func (tbr *tileBackgroundRender) Objects() []fyne.CanvasObject {
	tbr.mux.RLock()
	defer tbr.mux.RUnlock()
	log.Printf("objects return %d objects", len(tbr.overallContainer.Objects))
	// out, err := os.Create("./output.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var opt jpeg.Options
	// opt.Quality = 80

	// jpeg.Encode(out, newimg, &opt)
	return []fyne.CanvasObject{tbr.overallContainer}
}
func (tbr *tileBackgroundRender) Refresh() {
	log.Printf("refreshing")
	canvas.Refresh(tbr.tb)

}
