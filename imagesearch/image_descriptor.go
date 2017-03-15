package imagesearch

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"  //decode
	_ "image/jpeg" //decode
	_ "image/png"  //decode
	"io/ioutil"
)

type ImageDescriptor struct {
	Sha1  string
	HBins Bins
	SBins Bins
	VBins Bins
	QtyBins
}

/*
const (
	hBinsQty = 8 // 0 - 44, 45 - 89, ... 315 - 360
	sBinsQty = 12
	vBinsQty = 3
)
*/
func NewImageDescriptorFromFile(path string, qtyBins QtyBins) (*ImageDescriptor, error) {
	fileBytes, err := ioutil.ReadFile(path)
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

	return NewImageDescriptor(hash, picture, qtyBins), nil

}

func NewImageDescriptor(sha1 string, picture image.Image, qtyBins QtyBins) *ImageDescriptor {

	hBinsCounts := make([]uint, qtyBins.HBins)
	sBinsCounts := make([]uint, qtyBins.SBins)
	vBinsCounts := make([]uint, qtyBins.VBins)

	for y := 0; y < picture.Bounds().Max.Y; y++ {
		for x := 0; x < picture.Bounds().Max.X; x++ {
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
		}
	}

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

/*
// CalculateBinMatchQuality scores the match between the descriptor bins and the bins in the other descriptor.
func (descriptor *ImageDescriptor) CalculateBinMatchScore(otherDescriptor *ImageDescriptor) *MatchScore {
	var hScore, sScore, vScore float64

	for index, descriptorHBin := range descriptor.HBins {
		decimalDifference := math.Abs(descriptorHBin.CountAsDecimal-otherDescriptor.HBins[index].CountAsDecimal) + 1
		hScore = float64(1) / math.Pow(decimalDifference, 2)
		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 hue", descriptor.Sha1)
		}
	}

	for index, descriptorSBin := range descriptor.SBins {
		decimalDifference := math.Abs(descriptorSBin.CountAsDecimal-otherDescriptor.SBins[index].CountAsDecimal) + 1
		sScore += float64(1) / math.Pow(decimalDifference, 2)

		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 sat", descriptor.Sha1)
		}
	}

	for index, descriptorVBin := range descriptor.VBins {
		decimalDifference := math.Abs(descriptorVBin.CountAsDecimal-otherDescriptor.VBins[index].CountAsDecimal) + 1
		vScore += float64(1) / math.Pow(decimalDifference, 2)

		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 val", descriptor.Sha1)
		}
	}

	return NewMatchScore(hScore, sScore, vScore)
}
*/
