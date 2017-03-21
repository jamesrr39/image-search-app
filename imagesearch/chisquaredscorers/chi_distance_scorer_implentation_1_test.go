package chisquaredscorers

import (
	"image-search-app/imagesearch"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calcOneChiDistance(t *testing.T) {
	scorer := new(ChiDistanceSearchImpl1)

	distance1 := scorer.calcOneChiDistance(imagesearch.Bin{0}, imagesearch.Bin{0})
	assert.Equal(t, float64(0), distance1)

	distance2 := scorer.calcOneChiDistance(imagesearch.Bin{3}, imagesearch.Bin{4})
	assert.Equal(t, float64(0.125), distance2)

	distance3 := scorer.calcOneChiDistance(imagesearch.Bin{400}, imagesearch.Bin{4})
	assert.Equal(t, float64(387.2), distance3)
}
