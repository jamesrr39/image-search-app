package imagesearch

type ImageScorer interface {
	Score(seedImage, imageBeingScored *ImageDescriptor) MatchScore
}
