package imagesearch

import (
	"image/color"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewImageDescriptor(t *testing.T) {
	qtyBins := QtyBins{8, 12, 3}

	filePath := "IMG_20160616_130244.jpg"
	file, err := os.Open(filePath)
	assert.Nil(t, err)
	defer file.Close()

	descriptor, err := NewImageDescriptorFromFile(file, qtyBins)
	assert.Nil(t, err)

	for _, hBin := range descriptor.HBins {
		assert.False(t, math.IsInf(hBin.CountAsDecimal, 0))
	}

	for _, sBin := range descriptor.SBins {
		assert.False(t, math.IsInf(sBin.CountAsDecimal, 0))
	}

	for _, vBin := range descriptor.VBins {
		assert.False(t, math.IsInf(vBin.CountAsDecimal, 0))
	}
}

func Test_colorToRGBA(t *testing.T) {
	startColor := color.RGBA{23, 87, 122, 12}
	convertedColor := colorToRGBA(startColor)
	assert.Equal(t, startColor, convertedColor)
}
