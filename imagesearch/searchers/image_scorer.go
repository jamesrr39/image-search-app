package searchers

import (
	"image-search-app/imagesearch"
)

type DescriptorWithMatchScore struct {
	MatchScore *imagesearch.MatchScore
	Descriptor *imagesearch.PersistedImageDescriptor
}

type ImageScorer interface {
	Score(seedImage, imageBeingScored *imagesearch.ImageDescriptor) imagesearch.MatchScore
}
