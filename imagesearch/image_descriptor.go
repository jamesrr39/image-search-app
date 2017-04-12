package imagesearch

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"  //decode
	_ "image/jpeg" //decode
	_ "image/png"  //decode
	"io"
	"io/ioutil"

	"github.com/jamesrr39/semaphore"
)

type ImageDescriptor struct {
	Sha1  string
	HBins Bins
	SBins Bins
	VBins Bins
	QtyBins
}

func NewImageDescriptor(sha1 string, hBins, sBins, vBins Bins, qtyBins QtyBins) *ImageDescriptor {
	return &ImageDescriptor{sha1, hBins, sBins, vBins, qtyBins}
}

func NewImageDescriptorFromFile(reader io.Reader, qtyBins QtyBins, maxGoRoutines uint) (*ImageDescriptor, error) {
	fileBytes, err := ioutil.ReadAll(reader)
	if nil != err {
		return nil, err
	}

	hash, err := HashOfFile(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}

	picture, _, err := image.Decode(bytes.NewBuffer(fileBytes))
	if nil != err {
		return nil, err
	}

	return NewImageDescriptorFromPicture(hash, picture, qtyBins, maxGoRoutines), nil

}

func NewImageDescriptorFromPicture(sha1 string, picture image.Image, qtyBins QtyBins, maxGoRoutines uint) *ImageDescriptor {

	hBinsCounts := make([]uint, qtyBins.HBins)
	sBinsCounts := make([]uint, qtyBins.SBins)
	vBinsCounts := make([]uint, qtyBins.VBins)

	sema := semaphore.NewSemaphore(maxGoRoutines)

	for y := 0; y < picture.Bounds().Max.Y; y++ {
		for x := 0; x < picture.Bounds().Max.X; x++ {
			sema.Add()
			go func() {
				defer sema.Done()
				c := colorToRGBA(picture.At(x, y))
				hsvColor := NewHSVFromRGB(c)

				hbinIndex := uint((hsvColor.H / float64(360)) * float64(qtyBins.HBins))
				if hbinIndex == qtyBins.HBins {
					hbinIndex--
				}
				sbinIndex := uint((hsvColor.S / float64(1)) * float64(qtyBins.SBins))
				if sbinIndex == qtyBins.SBins {
					sbinIndex--
				}
				vbinIndex := uint((hsvColor.V / float64(1)) * float64(qtyBins.VBins))
				if vbinIndex == qtyBins.VBins {
					vbinIndex--
				}
				hBinsCounts[hbinIndex]++
				sBinsCounts[sbinIndex]++
				vBinsCounts[vbinIndex]++
			}()
		}
	}
	sema.Wait()

	return &ImageDescriptor{
		Sha1:  sha1,
		HBins: NewBins(hBinsCounts),
		SBins: NewBins(sBinsCounts),
		VBins: NewBins(vBinsCounts),
	}
}

func colorToRGBA(pixelColor color.Color) color.RGBA {
	r, g, b, a := pixelColor.RGBA()

	ratio8Bit32Bit := float64(255) / float64(65336)

	eightBitColour := color.RGBA{
		R: uint8(float64(r) * ratio8Bit32Bit),
		G: uint8(float64(g) * ratio8Bit32Bit),
		B: uint8(float64(b) * ratio8Bit32Bit),
		A: uint8(float64(a) * ratio8Bit32Bit),
	}
	return eightBitColour
}
