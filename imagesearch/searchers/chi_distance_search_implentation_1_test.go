package searchers

import (
	"image-search-app/imagesearch"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calcOneChiDistance(t *testing.T) {
	scorer := new(ChiDistanceSearchImpl1)

	seedBin := imagesearch.Bin{0}
	beingScoredBin := imagesearch.Bin{0}

	distance := scorer.calcOneChiDistance(seedBin, beingScoredBin)
	assert.False(t, math.IsNaN(distance))
}
