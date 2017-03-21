package imagesearchdal

import (
	"image-search-app/imagesearch"
)

type ImageSearchDAL interface {
	EnsureInCache(PersistedImageDescriptor) error
	Search(seedImageDescriptor *imagesearch.ImageDescriptor, scoringAlgorithm imagesearch.ImageScorer) ([]*imagesearch.DescriptorWithMatchScore, error)
}
