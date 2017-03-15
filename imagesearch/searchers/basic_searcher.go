package searchers

import (
	"image-search-app/imagesearch"
	"log"
	"math"
)

type BasicScorer struct{}

func (searcher *BasicScorer) Score(seedImage *imagesearch.ImageDescriptor, imageBeingScored *imagesearch.ImageDescriptor) imagesearch.MatchScore {
	var hScore, sScore, vScore float64

	for index, descriptorHBin := range seedImage.HBins {
		decimalDifference := math.Abs(descriptorHBin.CountAsDecimal-imageBeingScored.HBins[index].CountAsDecimal) + 1
		hScore = float64(1) / math.Pow(decimalDifference, 2)
		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 hue", seedImage.Sha1)
		}
	}

	for index, descriptorSBin := range seedImage.SBins {
		decimalDifference := math.Abs(descriptorSBin.CountAsDecimal-imageBeingScored.SBins[index].CountAsDecimal) + 1
		sScore += float64(1) / math.Pow(decimalDifference, 2)

		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 sat", seedImage.Sha1)
		}
	}

	for index, descriptorVBin := range seedImage.VBins {
		decimalDifference := math.Abs(descriptorVBin.CountAsDecimal-imageBeingScored.VBins[index].CountAsDecimal) + 1
		vScore += float64(1) / math.Pow(decimalDifference, 2)

		if decimalDifference == float64(0) {
			log.Printf("%s decimal difference was 0 val", seedImage.Sha1)
		}
	}

	return imagesearch.MatchScore((100 * hScore) + (10 * sScore) + vScore)
}
