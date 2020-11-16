package generator

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

const (
	_rows = 5
	_cols = 6
)

//identicon struct
type identicon struct {
	size    int
	cell    int
	bgcolor color.RGBA
	fgcolor color.RGBA
	image   *image.RGBA
	xmargin int
	ymargin int
}

func New(size int) *identicon {
	//cell side length = min(halfWidth/_cols,halfHeight/_rows)
	//half because the image is created for the half then mirrored
	cell := setCellSide(size)

	//make the image
	img := initImage(size)
	xmargin, ymargin := getMargins(size, cell)

	return &identicon{
		size: size,
		cell: cell,
		bgcolor: color.RGBA{
			R: 255, G: 255, B: 255, A: 0xff,
		},
		image:   img,
		xmargin: xmargin,
		ymargin: ymargin,
	}
}
func setCellSide(size int) int {
	//since the image will be vertically mirrored we choose cell size to be around half size of image
	halfSize := size / 2
	d1 := halfSize / (_cols / 2)
	d2 := halfSize / (_rows)
	if d1 < d2 {
		return d1
	}
	return d2
}

func initImage(size int) *image.RGBA {
	pointTop := image.Point{0, 0}
	pointBot := image.Point{size, size}
	return image.NewRGBA(image.Rectangle{pointTop, pointBot})
}

func getMargins(size, cell int) (xmargin, ymargin int) {
	//margin horizontal
	xmargin = size/2 - (cell*_cols)/2
	//margin vertical
	ymargin = size/2 - (cell*_rows)/2
	return
}

func (icon *identicon) renderCell(x, y int, colour color.RGBA) (availableX int) {
	//render a cell for given x,y and returns last x point available
	img := icon.image
	var i, j int
	for i = x; i-x <= icon.cell; i++ {
		for j = y; j-y <= icon.cell; j++ {
			img.SetRGBA(i, j, colour)
		}
	}
	return i
}

func (icon *identicon) mirrorHorizontally() {
	img := icon.image
	halfSize := icon.size / 2
	for x := icon.xmargin; x < halfSize; x++ {
		for y := icon.ymargin; y < icon.size; y++ {
			img.SetRGBA(icon.size-x, y, img.RGBAAt(x, y))
		}
	}
}

func (icon *identicon) render(hash []uint8, wg *sync.WaitGroup) {
	hashSliced := hash[3:]
	var x, y int
	x = icon.xmargin
	y = icon.ymargin
	presentByte := 0
	for i := 0; i < _rows; i++ {
		for j := 0; j < _cols/2; j++ {
			if hashSliced[presentByte]%2 == 0 {
				x = icon.renderCell(x, y, icon.fgcolor)
			} else {
				x = icon.renderCell(x, y, icon.bgcolor)
			}
			presentByte++
		}
		x = icon.xmargin
		y += icon.cell
	}
	if wg != nil {
		wg.Done()
	}
}
func (icon *identicon) renderBackground(wg *sync.WaitGroup) {
	img := icon.image
	//drawing top & bottom background
	for x := 0; x < icon.size; x++ {
		for y := 0; y < icon.ymargin; y++ {
			img.SetRGBA(x, y, icon.bgcolor)
			img.SetRGBA(x, icon.size-y, icon.bgcolor)
		}
	}
	//drawing left/right background
	for y := icon.ymargin; y <= icon.size-icon.ymargin; y++ {
		for x := 0; x <= icon.xmargin; x++ {
			img.SetRGBA(x, y, icon.bgcolor)
			img.SetRGBA(icon.size-x, y, icon.bgcolor)
		}
	}

	if wg != nil {
		wg.Done()
	}

}

func (icon *identicon) Create(hash []uint8, fileName string) interface{} {
	//Using first 3 hash bytes for color
	icon.fgcolor = color.RGBA{hash[0], hash[1], hash[2], 255}
	var waiter sync.WaitGroup
	waiter.Add(2)
	go icon.render(hash, &waiter)
	go icon.renderBackground(&waiter)
	waiter.Wait()
	icon.mirrorHorizontally()
	f, err := os.Create(fileName + ".png")
	if err != nil {
		return err
	} else {
		err = png.Encode(f, icon.image)
	}
	return err
}

func (icon *identicon) NoCreate(hash []uint8, fileName string) interface{} {
	//Using first 3 hash bytes for color
	icon.fgcolor = color.RGBA{hash[0], hash[1], hash[2], 255}

	icon.render(hash, nil)
	icon.renderBackground(nil)

	icon.mirrorHorizontally()
	f, err := os.Create(fileName + ".png")
	if err != nil {
		return err
	} else {
		err = png.Encode(f, icon.image)
	}
	return err
}
