package imagesearchfsdal

import (
	"image-search-app/imagesearch"

	"github.com/bradfitz/slice"
)

func (dal *ImageSearchFsDAL) Search(seedImageDescriptor *imagesearch.ImageDescriptor, scoringAlgorithm imagesearch.ImageScorer) []*DescriptorWithMatchScore {

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
