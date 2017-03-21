package imagesearch

type DescriptorWithMatchScore struct {
	MatchScore *MatchScore
	Descriptor *PersistedImageDescriptor
}

type ImageScorer interface {
	Score(seedImage, imageBeingScored *ImageDescriptor) MatchScore
}
