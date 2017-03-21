package chisquaredscorers

import (
	"image-search-app/imagesearch"
	"log"
	"math"
)

type ChiDistanceSearchImpl1 struct {
}

// (x-y)^2 /  x+y
func (search *ChiDistanceSearchImpl1) Score(seedImage, imageBeingScored *imagesearch.ImageDescriptor) imagesearch.MatchScore {
	var hScore, sScore, vScore float64
	for index, hBin := range seedImage.HBins {
		hScore += search.calcOneChiDistance(hBin, imageBeingScored.HBins[index])
	}

	for index, sBin := range seedImage.SBins {
		sScore += search.calcOneChiDistance(sBin, imageBeingScored.SBins[index])
	}

	for index, vBin := range seedImage.VBins {
		vScore += search.calcOneChiDistance(vBin, imageBeingScored.VBins[index])
	}

	if math.IsNaN(hScore) {
		log.Printf("NaN match. seed: %v, comparator %v\n", seedImage, imageBeingScored)
	}

	return imagesearch.MatchScore(hScore)
}

func (search *ChiDistanceSearchImpl1) calcOneChiDistance(seedBin, beingScoredBin imagesearch.Bin) float64 {

	topRow := math.Pow((seedBin.CountAsDecimal - beingScoredBin.CountAsDecimal), 2)
	bottomRow := seedBin.CountAsDecimal + beingScoredBin.CountAsDecimal + 1 // +1 to avoid divide by 0
	return topRow / bottomRow

}
