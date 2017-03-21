package imagesearchdal

import (
	"image-search-app/imagesearch"

	"image-search-app/imagesearch/searchers"

	"github.com/bradfitz/slice"
)

type DescriptorWithMatchScore struct {
	MatchScore imagesearch.MatchScore
	Descriptor *imagesearch.PersistedImageDescriptor
}

func (dal *ImageSearchDAL) Search(seedImageDescriptor *imagesearch.ImageDescriptor, scoringAlgorithm searchers.ImageScorer) []*DescriptorWithMatchScore {

	var descriptorsWithScore []*DescriptorWithMatchScore

	descriptorsInCache := dal.cache.GetAll()
	for _, descriptorInCache := range descriptorsInCache {
		matchScore := scoringAlgorithm.Score(seedImageDescriptor, descriptorInCache.ImageDescriptor)
		descriptorsWithScore = append(descriptorsWithScore, &DescriptorWithMatchScore{matchScore, descriptorInCache})
	}

	slice.Sort(descriptorsWithScore, func(i, j int) bool {
		a := descriptorsWithScore[i]
		b := descriptorsWithScore[j]
		aScore := float64(a.MatchScore)
		bScore := float64(b.MatchScore)

		if aScore == bScore {
			return a.Descriptor.Sha1 < b.Descriptor.Sha1 // for deterministicness
		}

		return aScore < bScore
	})

	return descriptorsWithScore
}
