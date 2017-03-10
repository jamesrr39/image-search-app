package imagesearchcache

import (
	"image-search-app/imagesearch"
	"log"

	"github.com/bradfitz/slice"
)

type DescriptorWithMatchScore struct {
	Score      float64
	Descriptor *PersistedImageDescriptor
}

func (cache *ImageSearchCache) Search(seedImageDescriptor *imagesearch.ImageDescriptor) []*DescriptorWithMatchScore {

	var descriptorsWithScore []*DescriptorWithMatchScore

	descriptorsInCache := cache.GetAll()
	for _, descriptorInCache := range descriptorsInCache {
		matchScore := descriptorInCache.CalculateBinMatchScore(seedImageDescriptor)
		descriptorsWithScore = append(descriptorsWithScore, &DescriptorWithMatchScore{matchScore, descriptorInCache})
	}

	slice.Sort(descriptorsWithScore, func(i, j int) bool {
		a := descriptorsWithScore[i]
		b := descriptorsWithScore[j]

		if a.Score == b.Score {
			return a.Descriptor.Sha1 > b.Descriptor.Sha1 // for deterministicness
		}

		return a.Score > b.Score
	})
	log.Println("after sorting")

	return descriptorsWithScore
}
